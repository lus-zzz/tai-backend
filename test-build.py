#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速测试编译脚本
编译当前平台版本并测试版本信息功能
"""

import os
import sys
import subprocess
import platform
from pathlib import Path


def run_command(cmd, cwd=None):
    """运行命令"""
    try:
        result = subprocess.run(
            cmd,
            cwd=cwd,
            capture_output=True,
            text=True,
            shell=True if platform.system() == 'Windows' else False
        )
        return result.returncode, result.stdout, result.stderr
    except Exception as e:
        return 1, "", str(e)


def main():
    script_dir = Path(__file__).parent.resolve()
    
    print("=" * 60)
    print("快速测试编译脚本")
    print("=" * 60)
    
    # 1. 编译当前平台版本
    print("\n[1/3] 编译当前平台版本...")
    build_cmd = [sys.executable, "build.py", "-p", "current", "-o", "test-build", "-c"]
    returncode, stdout, stderr = run_command(build_cmd, cwd=script_dir)
    
    if returncode != 0:
        print(f"❌ 编译失败: {stderr}")
        sys.exit(1)
    
    print("✅ 编译成功")
    
    # 2. 测试 --version 参数
    print("\n[2/3] 测试 --version 参数...")
    test_build_dir = script_dir / "test-build"
    
    if platform.system() == 'Windows':
        binary_path = test_build_dir / "chat-backend.exe"
    else:
        binary_path = test_build_dir / "chat-backend"
    
    if not binary_path.exists():
        print(f"❌ 找不到编译后的二进制文件: {binary_path}")
        sys.exit(1)
    
    version_cmd = [str(binary_path), "--version"]
    returncode, stdout, stderr = run_command(version_cmd)
    
    if returncode != 0:
        print(f"❌ 版本命令失败: {stderr}")
        sys.exit(1)
    
    print("✅ 版本信息输出:")
    print("-" * 40)
    print(stdout)
    print("-" * 40)
    
    # 3. 验证版本信息包含必要字段
    print("\n[3/3] 验证版本信息...")
    required_fields = ["Version:", "Build Time:", "Git Commit:"]
    missing_fields = []
    
    for field in required_fields:
        if field not in stdout:
            missing_fields.append(field)
    
    if missing_fields:
        print(f"❌ 缺少版本字段: {', '.join(missing_fields)}")
        sys.exit(1)
    
    print("✅ 版本信息验证通过")
    
    # 总结
    print("\n" + "=" * 60)
    print("✅ 所有测试通过!")
    print("=" * 60)
    print(f"\n编译文件位置: {binary_path}")
    print(f"文件大小: {binary_path.stat().st_size / (1024*1024):.2f} MB")
    print("\n可以使用以下命令测试运行:")
    if platform.system() == 'Windows':
        print(f"  cd test-build")
        print(f"  .\\chat-backend.exe")
    else:
        print(f"  cd test-build")
        print(f"  ./chat-backend")
    print("\n访问 API 文档: http://localhost:9090/swagger/index.html")
    print("查看版本信息: http://localhost:9090/api/v1/version")


if __name__ == '__main__':
    main()
