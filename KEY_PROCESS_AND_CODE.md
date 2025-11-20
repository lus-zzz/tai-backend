

## 核心功能流程

### 1. 文档分块向量化流程 (chunkAndVectorize)

该流程负责将 PDF 文档进行分块处理，并将分块结果向量化后存储到 Qdrant 向量数据库中。

#### 关键代码实现

```go
func chunkAndVectorize(docFile, collectionName string, cfg *config.Config) {
    fmt.Println("========== 分块和向量化模式 ==========")

    // 检查PDF文件是否存在
    if _, err := os.Stat(docFile); os.IsNotExist(err) {
        log.Fatalf("文件不存在: %s\n请检查文件路径是否正确，支持绝对路径或相对路径", docFile)
    }

    fmt.Printf("正在读取PDF文件: %s\n", docFile)
    // 初始化嵌入模型
    var ollamaEmbedder *embeddings.EmbedderImpl
    var err error


    ollamaEmbedderModel, err := ollama.New(ollama.WithServerURL(cfg.Embedding.Ollama.BaseURL), ollama.WithModel(cfg.Embedding.Ollama.Model))
    if err != nil {
        log.Fatal(err)
    }
    ret, err := ollamaEmbedderModel.CreateEmbedding(context.Background(), []string{"aaa"}) // 预热模型，避免首次调用延迟
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✓ 嵌入模型预热成功，返回示例向量长度: %d\n", len(ret))
    ollamaEmbedder, err = embeddings.NewEmbedder(ollamaEmbedderModel)
    if err != nil {
        log.Fatal(err)
    }
    

    // 使用Docling API分块文件
    doclingClient := chunk.NewDoclingClient(cfg.Docling.BaseURL, cfg.Docling.APIKey)

    // 读取PDF文件
    pdfData, err := os.ReadFile(docFile)
    if err != nil {
        log.Fatalf("无法读取PDF文件 %s: %v", docFile, err)
    }

    // 准备分块请求
    req := &chunk.HybridChunkerRequest{
        Files: [][]byte{pdfData},
        // 转换选项
        ConvertDoOCR:               false,
        ConvertImageExportMode:     "placeholder",
        ConvertPDFBackend:          "dlparse_v4",
        ConvertTableMode:           "accurate",
        ConvertPipeline:            "standard",
        ConvertAbortOnError:        false,
        ConvertDoCodeEnrichment:    false,
        ConvertDoFormulaEnrichment: false,
        // 分块选项
        ChunkingUseMarkdownTables: true,
        ChunkingIncludeRawText:    false,
        ChunkingTokenizer:  "/Volume2/test_work/models/BAAI/bge-m3",
        ChunkingMaxTokens:  512,  // 增加每块的最大token数，减少分块数量
        ChunkingMergePeers: true, // 合并相邻的小分块
    }

    // 调用同步分块API
    fmt.Println("正在分块文件...")
    startTime := time.Now()
    result, err := doclingClient.ChunkFilesWithHybridChunker(req)
    elapsedTime := time.Since(startTime)
    fmt.Printf("文件分块完成，耗时: %.2fs\n", elapsedTime.Seconds())

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("分块成功！共获得 %d 个分块\n", len(result.Chunks))



    // 将分块结果转换为schema.Document格式
    startTime = time.Now()
    docs := make([]schema.Document, 0, len(result.Chunks))
    for _, c := range result.Chunks {
        // 构建完整文本：包含元数据注释和分块内容
        pageStr := ""
        if c.PageNumbers != nil && len(*c.PageNumbers) > 0 {
            pageStr = fmt.Sprintf(", pages: %v", *c.PageNumbers)
        }
        fullText := fmt.Sprintf("<document_metadata> sourceDocument: %s published: %s %s </document_metadata>\n\n%s",
            c.Filename,
            time.Now().Format("2006/01/02 15:04:05"),
            pageStr,
            c.Text)

        // 生成description：取分块内容的前200个字符作为摘要
        description := c.Text
        if len(description) > 200 {
            description = description[:200] + "..."
        }

        // 构建严格的10字段元数据结构（符合AnythingLLM格式）
        metadata := map[string]interface{}{
            "url":                  "file://" + docFile,
            "title":                c.Filename,
            "docAuthor":            "",          // 可从PDF元数据提取
            "description":          description, // 从分块内容前200字符生成
            "docSource":            "pdf file uploaded by the user.",
            "chunkSource":          "localfile://" + docFile,
            "published":            time.Now().Format("2006/01/02 15:04:05"),
            "wordCount":            calculateWordCount(c.Text),
            "token_count_estimate": 0,
        }

        // 设置token计数
        if c.NumTokens != nil {
            metadata["token_count_estimate"] = *c.NumTokens
        }

        doc := schema.Document{
            PageContent: fullText,
            Metadata:    metadata,
        }
        docs = append(docs, doc)
    }
    elapsedTime = time.Since(startTime)
    fmt.Printf("将分块结果转换为schema.Document格式！耗时: %.2fs\n", elapsedTime.Seconds())

    // 将文档存储到Qdrant向量数据库
    fmt.Printf("正在将 %d 个文档向量化并存储到 Qdrant...\n", len(docs))
    startTime = time.Now()
    if useStorage(docs, collectionName, ollamaEmbedder, cfg) != nil {
        elapsedTime := time.Since(startTime)
        fmt.Printf("✓ 向量化存储完成！耗时: %.2fs\n", elapsedTime.Seconds())
    } else {
        fmt.Println("✗ 向量化存储失败！")
    }
}
```

#### 流程说明

1. **参数校验**: 检查指定的文件是否存在
2. **初始化嵌入模型**: 选择Ollama作为嵌入模型提供者
3. **文档分块**: 使用Docling API对PDF文件进行分块处理
5. **格式转换**: 将分块结果转换为 schema.Document 格式，并添加丰富元数据
6. **向量化存储**: 调用 useStorage 函数将文档存储到 Qdrant 向量数据库中

### 2. 问答交互流程 (chat)

该流程实现了基于已向量化的文档进行交互式问答的功能。

#### 关键代码实现

```go
func chat(collectionName string, cfg *config.Config) {
    fmt.Println("========== 问答模式 ==========")
    ctx := context.Background()

    // 初始化LLM模型
    startInit := time.Now()
    var llm llms.Model
    var err error


    llm, err = openai.New(
        openai.WithBaseURL(cfg.LLM.OpenAI.BaseURL),
        openai.WithToken(cfg.LLM.OpenAI.Token),
    )


    if err != nil {
        log.Fatal(err)
    }

    // 初始化嵌入模型
    ollamaEmbedderModel, err := ollama.New(ollama.WithServerURL(cfg.Embedding.Ollama.BaseURL), ollama.WithModel(cfg.Embedding.Ollama.Model))
    if err != nil {
        log.Fatal(err)
    }
    ollamaEmbedder, err = embeddings.NewEmbedder(ollamaEmbedderModel)
    if err != nil {
        log.Fatal(err)
    }
    

    fmt.Printf("✓ 模型初始化完成，耗时: %.2fs\n", time.Since(startInit).Seconds())

    // 连接到已有的Qdrant向量库
    startConnect := time.Now()
    fmt.Printf("正在连接到 Qdrant 集合: %s\n", collectionName)
    store := connectToQdrant(collectionName, ollamaEmbedder, cfg)
    if store == nil {
        log.Fatal("无法连接到 Qdrant")
    }
    fmt.Printf("✓ Qdrant 连接完成，耗时: %.2fs\n", time.Since(startConnect).Seconds())

    // 启动交互式问答循环
    fmt.Println("开始问答 (输入 'exit' 退出):")
    for {
        prompt, err := utils.GetUserInput("Q")
        if err != nil {
            log.Fatal(err)
        }

        if prompt == "exit" {
            fmt.Println("再见！")
            break
        }

        // 记录单次问答的总耗时
        roundStart := time.Now()

        // 使用向量搜索检索相关文档
        startRetrieval := time.Now()
        retrievalTime := time.Since(startRetrieval).Seconds()

        optionsVector := []vectorstores.Option{
            vectorstores.WithScoreThreshold(0.4), // 降低阈值，提高检索匹配率
        }

        retriever := vectorstores.ToRetriever(store, 10, optionsVector...)

        // 创建持久化到文件的 SQLite 记忆
        history := sqlite3.NewSqliteChatMessageHistory(
            sqlite3.WithDBAddress(cfg.ChatHistory.DBPath),
            sqlite3.WithSession(cfg.ChatHistory.Session),
            sqlite3.WithOverwrite(),
        )

        conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

        executor := chains.NewConversationalRetrievalQAFromLLM(
            llm,
            retriever,
            conversation,
        )

        options := []chains.ChainCallOption{
            chains.WithTemperature(0.8),
        }

        // 使用回调函数处理流式输出
        streamHandler := func(ctx context.Context, chunk []byte) error {
            fmt.Print(string(chunk))
            return nil
        }

        // 记录 LLM 处理时间
        startLLM := time.Now()
        _, err = chains.Run(ctx, executor, prompt,
            append(options, chains.WithStreamingFunc(streamHandler))...)
        llmTime := time.Since(startLLM).Seconds()

        if err != nil {
            fmt.Println("\nError running chains:", err)
            continue
        }

        roundTime := time.Since(roundStart).Seconds()
        fmt.Printf("\n  [LLM处理] 耗时: %.2fs\n", llmTime)
        fmt.Printf("  [本轮总耗时] %.2fs (检索: %.2fs, LLM: %.2fs)\n\n", roundTime, retrievalTime, llmTime)
    }
}
```

#### 流程说明

1. **模型初始化**: 根据配置初始化LLM（大语言模型）和嵌入模型
2. **连接向量数据库**: 使用 connectToQdrant 函数连接到现有的 Qdrant 集合
3. **交互式问答循环**: 循环接收用户输入的问题，当用户输入"exit"时退出程序
4. **检索增强生成**: 使用向量检索从 Qdrant 中查找相关文档，并结合 LLM 生成回答
5. **对话历史管理**: 使用 SQLite 数据库存储对话历史，确保应用重启后历史记录不会丢失
6. **问答执行**: 使用 NewConversationalRetrievalQAFromLLM 创建问答执行器，并支持流式输出
7. **性能统计**: 记录并显示每轮问答的处理时间，包括检索时间和LLM处理时间

## 辅助函数

### 向量存储 (useStorage)

```go
func useStorage(docs []schema.Document, collectionName string, embedder *embeddings.EmbedderImpl, cfg *config.Config) *qdrant1.Store {
    qdrantUrl, err := url.Parse(cfg.VectorStore.Qdrant.URL)
    if err != nil {
        log.Fatalf("failed parsing url: %s", err)
    }

    // Create collection if it doesn't exist
    err = ensureCollectionExists(qdrantUrl, collectionName, cfg.VectorStore.Qdrant.VectorSize)
    if err != nil {
        fmt.Println("Error ensuring collection exists:", err)
        return nil
    }

    store, err := qdrant1.New(
        qdrant1.WithURL(*qdrantUrl),
        qdrant1.WithAPIKey(cfg.VectorStore.Qdrant.APIKey),
        qdrant1.WithCollectionName(collectionName),
        qdrant1.WithEmbedder(embedder),
    )
    if err != nil {
        fmt.Println("Qdrant creation failed:", err)
        return nil
    }

    if len(docs) > 0 {
        // 分批处理文档，每批16个文档
        batchSize := 1000
        for i := 0; i < len(docs); i += batchSize {
            end := i + batchSize
            if end > len(docs) {
                end = len(docs)
            }

            batch := docs[i:end]
            _, err = store.AddDocuments(context.Background(), batch)
            if err != nil {
                fmt.Println("Error adding documents", err)
                return nil
            }

            // 如果有多批次，添加一个小延迟以减轻服务器压力
            if len(docs) > batchSize {
                fmt.Printf("Added batch %d-%d of documents\n", i+1, end)
            }
        }
    }

    return &store
}
```

### 连接向量数据库 (connectToQdrant)

```go
func connectToQdrant(collectionName string, embedder *embeddings.EmbedderImpl, cfg *config.Config) *qdrant1.Store {
    qdrantUrl, err := url.Parse(cfg.VectorStore.Qdrant.URL)
    if err != nil {
        log.Fatalf("failed parsing url: %s", err)
    }

    store, err := qdrant1.New(
        qdrant1.WithURL(*qdrantUrl),
        qdrant1.WithAPIKey(cfg.VectorStore.Qdrant.APIKey),
        qdrant1.WithCollectionName(collectionName),
        qdrant1.WithEmbedder(embedder),
    )
    if err != nil {
        fmt.Println("Qdrant connection failed:", err)
        return nil
    }

    return &store
}
```

## 总结

整个系统分为两个主要阶段：
1. **离线处理阶段**(chunkAndVectorize): 将文档分块、向量化并存储到向量数据库
2. **在线问答阶段**(chat): 基于已存储的向量数据进行交互式问答

系统充分利用了 LangchainGo 框架提供的各种组件，包括文档处理、向量存储、对话记忆和问答链等功能，构建了一个完整的检索增强生成(RAG)系统。