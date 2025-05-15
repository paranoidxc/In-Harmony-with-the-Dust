#!/bin/bash

# 检查是否安装了PlantUML
if ! command -v plantuml &> /dev/null; then
    echo "需要安装PlantUML来生成图表。"
    echo "可以通过以下方式安装："
    echo "  - 使用Homebrew: brew install plantuml"
    echo "  - 或者下载JAR文件: https://plantuml.com/download"
    exit 1
fi

# 生成类图
echo "正在生成类关系图..."
plantuml class_diagram.puml

echo "图表已生成: class_diagram.png"