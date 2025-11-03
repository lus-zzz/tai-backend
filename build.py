#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Chat Backend 编译脚本
支持 Windows 和 Linux 平台
"""

import os
import sys
import subprocess
import shutil
import argparse
import platform
import re
from datetime import datetime
from pathlib import Path


class Colors:
    """终端颜色输出"""
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

    @staticmethod
    def supports_color():
        """检查终端是否支持颜色"""
        plat = sys.platform
        supported_platform = plat != 'Pocket PC' and (plat != 'win32' or 'ANSICON' in os.environ)
        is_a_tty = hasattr(sys.stdout, 'isatty') and sys.stdout.isatty()
        return supported_platform and is_a_tty


# 如果不支持颜色，禁用所有颜色代码
if not Colors.supports_color() and platform.system() == 'Windows':
    # Windows 10+ 支持 ANSI 颜色
    try:
        import ctypes
        kernel32 = ctypes.windll.kernel32
        kernel32.SetConsoleMode(kernel32.GetStdHandle(-11), 7)
    except:
        for attr in dir(Colors):
            if not attr.startswith('_') and attr != 'supports_color':
                setattr(Colors, attr, '')


def print_info(msg):
    """打印信息"""
    print(f"{Colors.OKCYAN}[INFO]{Colors.ENDC} {msg}")


def print_success(msg):
    """打印成功信息"""
    print(f"{Colors.OKGREEN}[SUCCESS]{Colors.ENDC} {msg}")


def print_error(msg):
    """打印错误信息"""
    print(f"{Colors.FAIL}[ERROR]{Colors.ENDC} {msg}")


def print_warning(msg):
    """打印警告信息"""
    print(f"{Colors.WARNING}[WARNING]{Colors.ENDC} {msg}")


def run_command(cmd, cwd=None, env=None, capture_output=True):
    """运行命令并返回结果"""
    try:
        if capture_output:
            result = subprocess.run(
                cmd,
                cwd=cwd,
                env=env,
                capture_output=True,
                text=True,
                shell=True if platform.system() == 'Windows' else False
            )
            return result.returncode, result.stdout.strip(), result.stderr.strip()
        else:
            result = subprocess.run(cmd, cwd=cwd, env=env, shell=True if platform.system() == 'Windows' else False)
            return result.returncode, "", ""
    except Exception as e:
        return 1, "", str(e)


def get_git_commit():
    """获取 Git 提交哈希"""
    returncode, stdout, stderr = run_command(['git', 'rev-parse', '--short', 'HEAD'])
    if returncode == 0 and stdout:
        return stdout
    return "unknown"


def get_git_tag():
    """获取 Git 标签"""
    returncode, stdout, stderr = run_command(['git', 'describe', '--tags', '--abbrev=0'])
    if returncode == 0 and stdout:
        return stdout
    return ""


def get_auto_version():
    """获取自动版本号（基于提交计数，从0.0.0开始）"""
    # 获取总提交数
    returncode, stdout, stderr = run_command(['git', 'rev-list', '--count', 'HEAD'])
    if returncode == 0 and stdout:
        commit_count = int(stdout.strip())
        
        # 计算版本号组件 (从 0.0.0 开始)
        major = commit_count // 1000  # 每1000个提交进1
        minor = (commit_count % 1000) // 100  # 每100个提交进1
        patch = commit_count % 100  # 每个提交进1
        
        version = f"{major}.{minor}.{patch}"
        print_info(f"提交计数: {commit_count} -> 版本号: {version}")
        return version
    
    # 如果 Git 命令失败，使用时间戳
    timestamp = datetime.now().strftime('%Y%m%d%H%M')
    return f"0.0.{timestamp}"


def check_tag_exists(tag_name):
    """检查指定的标签是否已存在"""
    returncode, stdout, stderr = run_command(['git', 'tag', '-l', tag_name])
    if returncode == 0 and stdout.strip():
        return True
    return False


def get_git_branch():
    """获取 Git 分支"""
    returncode, stdout, stderr = run_command(['git', 'rev-parse', '--abbrev-ref', 'HEAD'])
    if returncode == 0 and stdout:
        return stdout
    return "unknown"


def get_git_status():
    """检查是否有未提交的更改"""
    returncode, stdout, stderr = run_command(['git', 'status', '--porcelain'])
    if returncode == 0 and stdout:
        return "dirty"
    return "clean"


def create_git_tag(tag_name, message=None):
    """创建 Git 标签"""
    if message is None:
        message = f"Release {tag_name}"
    
    # 检查标签是否已存在
    returncode, stdout, stderr = run_command(['git', 'tag', '-l', tag_name])
    if returncode == 0 and stdout.strip():
        print_warning(f"标签 {tag_name} 已存在")
        return False
    
    # 创建标签
    returncode, stdout, stderr = run_command(['git', 'tag', '-a', tag_name, '-m', message])
    if returncode != 0:
        print_error(f"创建标签失败: {stderr}")
        return False
    
    print_success(f"已创建标签: {tag_name}")
    return True


def push_git_tag(tag_name):
    """推送 Git 标签到远程仓库"""
    returncode, stdout, stderr = run_command(['git', 'push', 'origin', tag_name])
    if returncode != 0:
        print_error(f"推送标签失败: {stderr}")
        return False
    
    print_success(f"已推送标签到远程仓库: {tag_name}")
    return True


class Builder:
    """编译器类"""

    def __init__(self, args):
        self.args = args
        self.script_dir = Path(__file__).parent.resolve()
        self.chat_backend_dir = self.script_dir / 'chat-backend'  # chat-backend 子目录
        self.output_dir = self.script_dir / args.output  # 输出到项目根目录的 release 文件夹
        
        # 版本信息处理
        self.build_time = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        self.git_commit = get_git_commit()
        self.git_tag = get_git_tag()
        self.git_branch = get_git_branch()
        self.git_status = get_git_status()
        
        # 处理版本号
        if args.version == "auto":
            # 使用自动版本号
            self.version = get_auto_version()
            print_info(f"自动生成版本号: {self.version}")
        else:
            # 使用指定的版本号
            self.version = args.version
        
        # 如果状态是 dirty，在版本号后面加上标记
        if self.git_status == "dirty":
            self.version += "-dirty"
            print_warning(f"检测到未提交的更改，版本号标记为: {self.version}")
            print_info("提交更改后重新编译可创建正式标签")
        
        # 初始化 swagger 工具
        self.swagger_available = False
        if not args.skip_swagger:
            self.init_swagger_tool()

    def init_swagger_tool(self):
        """初始化 swagger 工具"""
        print_info("检查 swagger 工具...")
        returncode, stdout, stderr = run_command(['swagger', 'version'])
        
        if returncode != 0:
            print_warning("未找到 swagger 工具，尝试安装...")
            returncode, stdout, stderr = run_command(
                ['go', 'install', 'github.com/go-swagger/go-swagger/cmd/swagger@latest']
            )
            if returncode != 0:
                print_error(f"swagger 安装失败: {stderr}")
                print_warning("Swagger 功能将不可用")
                self.swagger_available = False
                return
            
            # 再次检查
            returncode, stdout, stderr = run_command(['swagger', 'version'])
            if returncode != 0:
                print_warning("swagger 安装后仍无法使用，请确保 $GOPATH/bin 在 PATH 中")
                self.swagger_available = False
                return
            
            print_success(f"swagger 工具安装成功: {stdout}")
            self.swagger_available = True
        else:
            print_success(f"swagger 工具已就绪: {stdout}")
            self.swagger_available = True

    def handle_git_tagging(self):
        """处理 Git 标签"""
        if self.args.no_tag:
            return True
        
        # 检查是否有未提交的更改
        if self.git_status == "dirty":
            print_warning("存在未提交的更改，跳过标签创建和推送")
            print_info("编译完成，但未创建Git标签")
            print_info("提交更改后重新编译可创建正式标签并推送")
            return True  # 返回True，不阻止编译完成
        
        # 检查是否在主分支或发布分支
        allowed_branches = ['main', 'master', 'release', 'develop']
        if self.git_branch not in allowed_branches:
            print_warning(f"当前分支 '{self.git_branch}' 不在推荐的发布分支列表中: {allowed_branches}")
            if not self.args.force_tag:
                print_info("使用 --force-tag 参数可强制在当前分支创建标签")
                return False
        
        tag_name = f"v{self.version}"
        
        # 检查标签是否已存在
        if check_tag_exists(tag_name):
            print_info(f"标签 {tag_name} 已存在，跳过创建和推送")
            print_info("代码没有变化，无需重复发布")
            return True
        
        # 创建标签
        if create_git_tag(tag_name, f"Release {self.version}"):
            # 默认推送标签到远程，除非明确禁用
            if not self.args.no_push:
                if not push_git_tag(tag_name):
                    print_warning("标签创建成功但推送失败")
                    return False
            else:
                print_info("标签已创建，但未推送到远程仓库（使用了 --no-push 参数）")
            return True
        else:
            print_error("标签创建失败")
            return False

    def print_build_info(self):
        """打印编译信息"""
        print(f"\n{Colors.BOLD}{'='*60}{Colors.ENDC}")
        print(f"{Colors.BOLD}Chat Backend 编译脚本{Colors.ENDC}")
        print(f"{Colors.BOLD}{'='*60}{Colors.ENDC}")
        print_info(f"版本: {Colors.BOLD}{self.version}{Colors.ENDC}")
        print_info(f"构建时间: {self.build_time}")
        print_info(f"Git Commit: {self.git_commit}")
        print_info(f"Git Branch: {self.git_branch}")
        print_info(f"Git Tag: {self.git_tag or 'N/A'}")
        print_info(f"Git Status: {self.git_status}")
        print_info(f"输出目录: {self.output_dir}")
        print(f"{Colors.BOLD}{'='*60}{Colors.ENDC}\n")

    def clean(self):
        """清理输出目录"""
        if self.output_dir.exists():
            print_info(f"清理输出目录: {self.output_dir}")
            shutil.rmtree(self.output_dir)
            print_success("清理完成")

    def prepare_output_dir(self):
        """准备输出目录"""
        if not self.output_dir.exists():
            print_info(f"创建输出目录: {self.output_dir}")
            self.output_dir.mkdir(parents=True, exist_ok=True)

    def check_go_version(self):
        """检查 Go 版本"""
        print_info("检查 Go 环境...")
        returncode, stdout, stderr = run_command(['go', 'version'])
        if returncode != 0:
            print_error("未找到 Go 环境，请先安装 Go")
            sys.exit(1)
        print_success(f"Go 环境: {stdout}")

    def download_dependencies(self):
        """下载依赖"""
        print_info("下载依赖...")
        returncode, stdout, stderr = run_command(['go', 'mod', 'download'], cwd=self.chat_backend_dir)
        if returncode != 0:
            print_error(f"依赖下载失败: {stderr}")
            sys.exit(1)
        print_success("依赖下载完成")

    def generate_swagger_docs(self):
        """生成 Swagger 文档"""
        if not self.swagger_available:
            print_warning("swagger 工具不可用，跳过 Swagger 文档生成")
            return False
        
        print_info("生成 Swagger 文档...")
        
        # 运行 swagger generate spec
        output_file = self.chat_backend_dir / 'docs' / 'swagger.json'
        cmd = ['swagger', 'generate', 'spec', '-o', str(output_file), '--scan-models']
        
        returncode, stdout, stderr = run_command(cmd, cwd=self.chat_backend_dir)
        if returncode != 0:
            print_error(f"Swagger 文档生成失败: {stderr}")
            print_warning("跳过 Swagger 文档生成，继续编译...")
            return False
        
        print_success("Swagger 文档生成完成")
        return True

    def build_binary(self, goos, goarch, output_name):
        """编译二进制文件"""
        print_info(f"编译 {goos}/{goarch} 版本...")
        
        # 构建环境变量
        env = os.environ.copy()
        env['GOOS'] = goos
        env['GOARCH'] = goarch
        env['CGO_ENABLED'] = '0'
        
        # 构建 ldflags
        ldflags = [
            "-s",  # 去除符号表
            "-w",  # 去除调试信息
            f"-X 'main.Version={self.version}'",
            f"-X 'main.BuildTime={self.build_time}'",
            f"-X 'main.GitCommit={self.git_commit}'",
            f"-X 'main.GitBranch={self.git_branch}'",
            f"-X 'main.GitTag={self.git_tag}'",
        ]
        ldflags_str = ' '.join(ldflags)
        
        # 输出文件路径
        output_path = self.output_dir / output_name
        
        # 编译命令
        cmd = [
            'go', 'build',
            '-ldflags', ldflags_str,
            '-trimpath',
            '-o', str(output_path),
            '.'
        ]
        
        returncode, stdout, stderr = run_command(cmd, cwd=self.chat_backend_dir, env=env)
        
        if returncode != 0:
            print_error(f"{goos}/{goarch} 编译失败: {stderr}")
            return False
        
        # 获取文件大小
        file_size = output_path.stat().st_size / (1024 * 1024)  # MB
        print_success(f"{goos}/{goarch} 编译完成: {output_name} ({file_size:.2f} MB)")
        return True

    def build_all_platforms(self):
        """编译所有平台"""
        platforms = []
        
        if self.args.platform == 'all':
            platforms = [
                ('windows', 'amd64', 'chat-backend-windows-amd64.exe'),
                ('linux', 'amd64', 'chat-backend-linux-amd64'),
                ('darwin', 'amd64', 'chat-backend-darwin-amd64'),
                ('darwin', 'arm64', 'chat-backend-darwin-arm64'),
            ]
        elif self.args.platform == 'windows':
            platforms = [('windows', 'amd64', 'chat-backend-windows-amd64.exe')]
        elif self.args.platform == 'linux':
            platforms = [('linux', 'amd64', 'chat-backend-linux-amd64')]
        elif self.args.platform == 'darwin':
            platforms = [
                ('darwin', 'amd64', 'chat-backend-darwin-amd64'),
                ('darwin', 'arm64', 'chat-backend-darwin-arm64'),
            ]
        elif self.args.platform == 'current':
            current_os = platform.system().lower()
            if current_os == 'windows':
                platforms = [('windows', 'amd64', 'chat-backend.exe')]
            elif current_os == 'linux':
                platforms = [('linux', 'amd64', 'chat-backend')]
            elif current_os == 'darwin':
                arch = platform.machine().lower()
                if 'arm' in arch or 'aarch64' in arch:
                    platforms = [('darwin', 'arm64', 'chat-backend')]
                else:
                    platforms = [('darwin', 'amd64', 'chat-backend')]
        
        success_count = 0
        for goos, goarch, output_name in platforms:
            if self.build_binary(goos, goarch, output_name):
                success_count += 1
        
        return success_count, len(platforms)

    def build(self):
        """执行编译"""
        try:
            # 打印编译信息
            self.print_build_info()
            
            # 清理
            if self.args.clean:
                self.clean()
            
            # 准备输出目录
            self.prepare_output_dir()
            
            # 检查 Go 环境
            self.check_go_version()
            
            # 下载依赖
            if not self.args.skip_deps:
                self.download_dependencies()
            
            # 生成 Swagger 文档
            if not self.args.skip_swagger:
                self.generate_swagger_docs()
            
            # 编译二进制
            success_count, total_count = self.build_all_platforms()
            
            if success_count == 0:
                print_error("所有平台编译失败")
                sys.exit(1)
            
            # 编译成功后处理 Git 标签（默认启用）
            print_info("开始处理 Git 标签...")
            if not self.handle_git_tagging():
                print_warning("Git 标签处理失败，但编译已完成")
            
            # 显示结果
            print(f"\n{Colors.BOLD}{'='*60}{Colors.ENDC}")
            print_success(f"编译完成! 成功 {success_count}/{total_count}")
            print_info(f"输出目录: {self.output_dir}")
            print(f"\n{Colors.BOLD}编译的文件:{Colors.ENDC}")
            
            for item in sorted(self.output_dir.iterdir()):
                if item.is_file() and 'chat-backend' in item.name and item.suffix in ['.exe', '']:
                    file_size = item.stat().st_size / (1024 * 1024)
                    print(f"  {Colors.WARNING}→{Colors.ENDC} {item.name} ({file_size:.2f} MB)")
            
            print(f"\n{Colors.BOLD}注意:{Colors.ENDC}")
            print_info("所有资源文件(配置、静态文件)已嵌入到二进制中")
            print_info("只需分发二进制文件即可运行")
            print(f"{Colors.BOLD}{'='*60}{Colors.ENDC}\n")
            
        except KeyboardInterrupt:
            print_warning("\n编译已取消")
            sys.exit(1)
        except Exception as e:
            print_error(f"编译失败: {e}")
            import traceback
            traceback.print_exc()
            sys.exit(1)


def main():
    parser = argparse.ArgumentParser(description='Chat Backend 编译脚本')
    parser.add_argument('-v', '--version', default='auto', 
                        help='版本号 (默认: auto - 自动生成版本号)')
    parser.add_argument('-o', '--output', default='release', help='输出目录 (默认: release)')
    parser.add_argument('-p', '--platform', choices=['all', 'windows', 'linux', 'darwin', 'current'], 
                        default='current', help='编译平台 (默认: current)')
    parser.add_argument('-c', '--clean', action='store_true', help='清理之前的构建')
    parser.add_argument('--skip-deps', action='store_true', help='跳过依赖下载')
    parser.add_argument('--skip-swagger', action='store_true', help='跳过 Swagger 文档生成')
    
    # Git 标签相关参数
    parser.add_argument('--no-tag', action='store_true', help='跳过自动创建和推送 Git 标签')
    parser.add_argument('--no-push', action='store_true', help='创建标签但不推送到远程仓库')
    parser.add_argument('--force-tag', action='store_true', help='强制在当前分支创建标签，忽略分支检查')
    
    args = parser.parse_args()
    
    # 参数验证
    if args.no_push and args.no_tag:
        print_error("--no-push 和 --no-tag 不能同时使用")
        sys.exit(1)
    
    builder = Builder(args)
    builder.build()


if __name__ == '__main__':
    main()
