# 构建和部署文档脚本

# 检查是否安装了 mkdocs 和 material 主题
check_installation() {
    if ! command -v mkdocs &> /dev/null; then
        echo "错误：mkdocs 未安装"
        echo "请运行：pip install mkdocs mkdocs-material"
        exit 1
    fi
}

# 构建文档
build_docs() {
    echo "正在构建文档..."
    mkdocs build --clean
    echo "文档构建完成！输出目录：site/"
}

# 本地预览
serve_docs() {
    echo "启动本地预览服务器..."
    echo "访问：http://127.0.0.1:8000"
    mkdocs serve
}

# 部署到 GitHub Pages
deploy_github_pages() {
    echo "部署到 GitHub Pages..."
    
    # 检查是否安装了 mike
    if ! command -v mike &> /dev/null; then
        echo "安装 mike..."
        pip install mike
    fi
    
    # 部署
    mike deploy --push --update-aliases main latest
    echo "部署完成！"
}

# 显示帮助
show_help() {
    echo "用法：./deploy.sh [命令]"
    echo ""
    echo "命令:"
    echo "  build    构建文档"
    echo "  serve    本地预览"
    echo "  deploy   部署到 GitHub Pages"
    echo "  help     显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  ./deploy.sh build   # 构建文档"
    echo "  ./deploy.sh serve   # 本地预览"
    echo "  ./deploy.sh deploy  # 部署到 GitHub Pages"
}

# 主逻辑
main() {
    check_installation
    
    case "${1:-help}" in
        build)
            build_docs
            ;;
        serve)
            serve_docs
            ;;
        deploy)
            deploy_github_pages
            ;;
        help|*)
            show_help
            ;;
    esac
}

main "$@"
