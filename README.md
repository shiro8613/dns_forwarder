# dns_forwarder

## これなに？

普通はdns_recursorと呼ばれるもので以下のような動きをします。

example.com -> 172.16.0.1:53

それ以外 -> 8.8.8.8

v6にも対応しているほか、複数のサーバーを登録することで、クエリに失敗した際に別な鯖へ問い合わせ可能です。

## インストール方法

Ubuntu
```shell
sudo mkdir -p /srv/dns_forwarder
sudo curl -L -o /srv/dns_forwarder/dns_forwarder "https://github.com/shiro8613/dns_forwarder/releases/latest/download/dns_forwarder_linux_$([[ "$(uname -m)" == "x86_64" ]] && echo "amd64" || echo "arm64")"
sudo chmod u+x /srv/dns_forwarder/dns_forwarder
sudo curl -L -o /etc/systemd/system/dns_forwarder.service https://raw.githubusercontent.com/shiro8613/dns_forwarder/main/dns_forwarder.service
sudo curl -L -o /srv/dns_forwarder/config.yml https://raw.githubusercontent.com/shiro8613/dns_forwarder/main/config.example.yml
```

起動
```
sudo systemctl enable --now dns_forwarder
sudo systemctl start dns_forwarder
```

### コンフィグファイルについて

コンフィグファイルは`config.example.yml`を参照して書いていただきたいのですが、補足があるため記述します。

serversの指定方法は
```
servers:
    ドメイン名:
        to: 
         - アドレス:ポート
        priority: 1
```
のように記述します。

priorityとはドメイン名がかぶった際の優先度です。

bbb.aaa.aとaaa.aをそれぞれ別な鯖に問い合わせたいときのために利用します。

普段は`1`でいいと思います。


Config

| キー       | 型                | 必須 | 補足                                                                                               | 
| ---------- | ----------------- | ---- | -------------------------------------------------------------------------------------------------- | 
| bind4      | string            | Yes  | v4のアドレスを`アドレス:ポート`で指定します。<br>ex) `0.0.0.0:53`                                      | 
| bind6      | string            | No   | v6のアドレスを`[アドレス]:ポート`でしていします。<br>""にするとv6を使用しなくなります。<br>ex) `[::]:53` | 
| cache_time | int               | Yes  | キャッシュの残留期間を秒で指定します。                                                             | 
| forwards   | []string          | Yes  | 登録されていないドメインの再問合せさきです。ex) 8.8.8.8                                            | 
| servers    | map[string]Server | Yes  | Server型は下に記述                                                                                 | 

Server

| キー     | 型       | 必須 | 補足                                                                                                                       | 
| -------- | -------- | ---- | -------------------------------------------------------------------------------------------------------------------------- | 
| to       | []string | Yes  | ドメインの問い合わせ先アドレスです。複数設定可能です。アドレスは`アドレス:ポート`でして指定します。<br>ex) `172.16.0.2:53` | 
| priority | int      | Yes  | ドメインがかぶっている場合の優先度をしていします。                                                                         | 

## ポートについて
説明では`アドレス:ポート`と書いていますが、`forwards`と`server`の`to`に関してはポートがない場合は`53`が適応されます。53を利用する場合は省略可能です。
