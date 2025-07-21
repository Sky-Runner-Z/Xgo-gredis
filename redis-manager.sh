#!/bin/bash
case "$1" in
    start)
        sudo systemctl start redis-server
        echo "Redis已启动"
        redis-cli ping
        ;;
    stop)
        sudo systemctl stop redis-server
        echo "Redis已停止"
        ;;
    status)
        sudo systemctl status redis-server
        ;;
    restart)
        sudo systemctl restart redis-server
        echo "Redis已重启"
        redis-cli ping
        ;;
    *)
        echo "用法: $0 {start|stop|status|restart}"
        exit 1
        ;;
esac
