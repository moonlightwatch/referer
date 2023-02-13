# Referer 防盗链

通过 referer 进行防盗链保护

## 配置举例

```yaml
# Static configuration

experimental:
  plugins:
    referer:
        moduleName: github.com/moonlightwatch/referer
        version: v0.1.4
```

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    my-plugin:
      plugin:
        referer:
          Type: white
          EmptyReferer: true
          Domains:
            - "*.baidu.com"
            - "google.com"
```

## 说明



Type 字段，可以选填：`white` 或者 `black`，分别表示白名单模式和黑名单模式。

白名单模式下：

1. 匹配 `Domains` 所载域名的 `referer` 参数，可以进行访问。
2. 若 `EmptyReferer` 为 `true` 则允许 `referer` 为空的请求进行访问。否则，不允许 `referer` 为空的请求进行访问。

黑名单模式下：

1. 匹配 `Domains` 所载域名的 `referer` 参数，拒绝访问。
2. 若 `EmptyReferer` 为 `true` 则拒绝 `referer` 为空的请求进行访问。否则，允许 `referer` 为空的请求进行访问。


Domains 字段，是一个列表，可以填写域名或者以 `*` 为通配符的子域名。