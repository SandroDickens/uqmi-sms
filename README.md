# uqmi-sms

Use UQMI to send, receive, and delete SMS

## 1. Build

```shell
go get uqmi-sms
```

```shell
go build
```

## 2. Read SMS

### 2.1 Read SMS by ID

```shell
uqmi-sms -read -id <ID>
```

### 2.2 Read all SMS

```shell
uqmi-sms -read
```

## 3. Delete SMS

### 3.1 Delete SMS by ID

```shell
uqmi-sms -delete -id <ID>
```

### 3.2 Delete all SMS

```shell
uqmi-sms -delete
```

## 4. Send SMS

```shell
uqmi-sms -send -target <target_phone_number> -text <message>
```

eg:

```shell
uqmi-sms -send -target 10010 -text "hello world!"
```
