# 常见编解码工具

对照Yakit Codec  编了一个WEB版编解码工具，支持多种编码和解码操作。

## 功能特性

### 编码功能 (15种)

- **Base64编码**: 标准Base64编码
- **HTML实体编码(强制)**: 将所有字符转换为HTML实体
- **HTML实体编码(强制十六进制)**: 使用十六进制格式的HTML实体
- **HTML实体编码(特殊字符)**: 只编码HTML特殊字符 (&lt;, &gt;, &amp;, &quot;, &#39;)
- **URL编码(强制)**: 将所有字符编码为URL格式 (%XX)
- **URL编码(特殊字符)**: 只编码需要URL编码的字符
- **URL路径编码(特殊字符)**: 用于URL路径的编码
- **双重URL编码**: 两次URL编码
- **十六进制编码**: 将文本转换为十六进制
- **Unicode中文编码**: 将非ASCII字符转换为Unicode转义序列
- **MD5编码**: MD5哈希
- **SM3编码**: SM3哈希
- **SHA1编码**: SHA1哈希
- **SHA-256编码**: SHA256哈希
- **SHA-512编码**: SHA512哈希

### 解码功能 (7种)

- **Base64解码**: Base64解码
- **HTML解码**: HTML实体解码
- **URL解码**: URL解码
- **URL路径解码**: URL路径解码
- **双重URL解码**: 两次URL解码
- **十六进制解码**: 十六进制解码
- **Unicode中文解码**: Unicode转义序列解码

## 使用方法

1. 启动服务器：

   ```bash
   go run main.go
   ```

2. 在浏览器中访问：http://localhost:18080

3. 选择相应的编码或解码功能，输入内容并处理

<img width="1324" height="821" alt="image" src="https://github.com/user-attachments/assets/04d46708-19e0-4985-93f4-c83c0c75773a" />

SM3有问题，懒得搞了，自用，如有bug自己想办法吧
