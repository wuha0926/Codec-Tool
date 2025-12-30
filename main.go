package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

type CodecTool struct{}

func main() {
	tool := &CodecTool{}

	// 静态文件处理
	http.HandleFunc("/", tool.homeHandler)
	http.HandleFunc("/encode", tool.encodeHandler)
	http.HandleFunc("/decode", tool.decodeHandler)

	fmt.Println("Web编解码工具启动，访问 http://localhost:18080")
	http.ListenAndServe(":18080", nil)
}

func (tool *CodecTool) homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>编解码工具</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }
        .container { display: flex; gap: 20px; }
        .section { flex: 1; background: #f5f5f5; padding: 20px; border-radius: 8px; }
        .textarea { width: 100%; height: 200px; padding: 10px; border: 1px solid #ddd; border-radius: 4px; font-family: monospace; }
        .select { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; margin: 10px 0; }
        .button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; margin: 5px; }
        .button:hover { background: #0056b3; }
        .result { background: #e9ecef; padding: 10px; border-radius: 4px; margin-top: 10px; word-break: break-all; }
        h1 { text-align: center; color: #333; }
        h2 { color: #666; border-bottom: 2px solid #007bff; padding-bottom: 5px; }
        .function-list { font-size: 14px; line-height: 1.6; }
    </style>
</head>
<body>
    <h1>🔧 编解码工具</h1>
    
    <div class="container">
        <div class="section">
            <h2>编码工具</h2>
            <form action="/encode" method="post">
                <label>输入内容:</label>
                <textarea name="input" class="textarea" placeholder="请输入要编码的内容..."></textarea>
                
                <label>选择编码方式:</label>
                <select name="method" class="select">
                    <option value="base64">Base64编码</option>
                    <option value="html_force">HTML实体编码(强制)</option>
                    <option value="html_force_hex">HTML实体编码(强制十六进制)</option>
                    <option value="html_special">HTML实体编码(特殊字符)</option>
                    <option value="url_force">URL编码(强制)</option>
                    <option value="url_special">URL编码(特殊字符)</option>
                    <option value="url_path_special">URL路径编码(特殊字符)</option>
                    <option value="double_url">双重URL编码</option>
                    <option value="hex">十六进制编码</option>
                    <option value="unicode">Unicode中文编码</option>
                    <option value="md5">MD5编码</option>
                    <option value="sm3">SM3编码</option>
                    <option value="sha1">SHA1编码</option>
                    <option value="sha256">SHA-256编码</option>
                    <option value="sha512">SHA-512编码</option>
                </select>
                
                <button type="submit" class="button">编码</button>
            </form>
        </div>
        
        <div class="section">
            <h2>解码工具</h2>
            <form action="/decode" method="post">
                <label>输入内容:</label>
                <textarea name="input" class="textarea" placeholder="请输入要解码的内容..."></textarea>
                
                <label>选择解码方式:</label>
                <select name="method" class="select">
                    <option value="base64">Base64解码</option>
                    <option value="html">HTML解码</option>
                    <option value="url">URL解码</option>
                    <option value="url_path">URL路径解码</option>
                    <option value="double_url">双重URL解码</option>
                    <option value="hex">十六进制解码</option>
                    <option value="unicode">Unicode中文解码</option>
                </select>
                
                <button type="submit" class="button">解码</button>
            </form>
        </div>
    </div>
    
    <div class="section" style="margin-top: 20px;">
        <h2>功能说明</h2>
        <div class="function-list">
            <h3>编码功能 (15种):</h3>
            <ul>
                <li><strong>Base64编码</strong>: 标准的Base64编码</li>
                <li><strong>HTML实体编码</strong>: 将所有字符转换为HTML实体</li>
                <li><strong>HTML实体编码(十六进制)</strong>: 使用十六进制格式的HTML实体</li>
                <li><strong>HTML实体编码(特殊字符)</strong>: 只编码HTML特殊字符 (&lt;, &gt;, &amp;, &quot;, &#39;)</li>
                <li><strong>URL编码</strong>: 强制URL编码所有字符</li>
                <li><strong>URL编码(特殊字符)</strong>: 只编码需要URL编码的字符</li>
                <li><strong>URL路径编码</strong>: 用于URL路径的编码</li>
                <li><strong>双重URL编码</strong>: 两次URL编码</li>
                <li><strong>十六进制编码</strong>: 将文本转换为十六进制</li>
                <li><strong>Unicode中文编码</strong>: 将非ASCII字符转换为Unicode转义序列</li>
                <li><strong>MD5编码</strong>: MD5哈希</li>
                <li><strong>SM3编码</strong>: SM3哈希(使用MD5替代)</li>
                <li><strong>SHA1编码</strong>: SHA1哈希</li>
                <li><strong>SHA-256编码</strong>: SHA256哈希</li>
                <li><strong>SHA-512编码</strong>: SHA512哈希</li>
            </ul>
            
            <h3>解码功能 (7种):</h3>
            <ul>
                <li><strong>Base64解码</strong>: Base64解码</li>
                <li><strong>HTML解码</strong>: HTML实体解码</li>
                <li><strong>URL解码</strong>: URL解码</li>
                <li><strong>URL路径解码</strong>: URL路径解码</li>
                <li><strong>双重URL解码</strong>: 两次URL解码</li>
                <li><strong>十六进制解码</strong>: 十六进制解码</li>
                <li><strong>Unicode解码</strong>: Unicode转义序列解码</li>
            </ul>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (tool *CodecTool) encodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	input := r.FormValue("input")
	method := r.FormValue("method")

	var result string
	var err error

	switch method {
	case "base64":
		result, err = tool.encodeBase64(input)
	case "html_force":
		result, err = tool.encodeHTMLForce(input)
	case "html_force_hex":
		result, err = tool.encodeHTMLForceHex(input)
	case "html_special":
		result, err = tool.encodeHTMLSpecial(input)
	case "url_force":
		result, err = tool.encodeURLForce(input)
	case "url_special":
		result, err = tool.encodeURLSpecial(input)
	case "url_path_special":
		result, err = tool.encodeURLPathSpecial(input)
	case "double_url":
		result, err = tool.encodeDoubleURL(input)
	case "hex":
		result, err = tool.encodeHex(input)
	case "unicode":
		result, err = tool.encodeUnicode(input)
	case "md5":
		result, err = tool.encodeMD5(input)
	case "sm3":
		result, err = tool.encodeSM3(input)
	case "sha1":
		result, err = tool.encodeSHA1(input)
	case "sha256":
		result, err = tool.encodeSHA256(input)
	case "sha512":
		result, err = tool.encodeSHA512(input)
	default:
		err = fmt.Errorf("未知的编码方法")
	}

	tool.showResult(w, "编码", input, method, result, err)
}

func (tool *CodecTool) decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	input := r.FormValue("input")
	method := r.FormValue("method")

	var result string
	var err error

	switch method {
	case "base64":
		result, err = tool.decodeBase64(input)
	case "html":
		result, err = tool.decodeHTML(input)
	case "url":
		result, err = tool.decodeURL(input)
	case "url_path":
		result, err = tool.decodeURLPath(input)
	case "double_url":
		result, err = tool.decodeDoubleURL(input)
	case "hex":
		result, err = tool.decodeHex(input)
	case "unicode":
		result, err = tool.decodeUnicode(input)
	default:
		err = fmt.Errorf("未知的解码方法")
	}

	tool.showResult(w, "解码", input, method, result, err)
}

func (tool *CodecTool) showResult(w http.ResponseWriter, opType, input, method, result string, err error) {
	html := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>结果 - 编解码工具</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1000px; margin: 0 auto; padding: 20px; background: #f8f9fa; }
        .container { background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .section { margin: 20px 0; }
        .label { font-weight: bold; color: #495057; margin-bottom: 5px; display: block; }
        .content { background: #f8f9fa; padding: 15px; border-radius: 5px; border-left: 4px solid #007bff; word-break: break-all; font-family: monospace; }
        .result { background: #d4edda; border-left-color: #28a745; }
        .error { background: #f8d7da; border-left-color: #dc3545; }
        .button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; text-decoration: none; display: inline-block; margin: 10px 5px; }
        .button:hover { background: #0056b3; }
        .button-secondary { background: #6c757d; }
        .button-secondary:hover { background: #545b62; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ %s完成</h1>
            <p>方法: %s</p>
        </div>
        
        <div class="section">
            <span class="label">原始输入:</span>
            <div class="content">%s</div>
        </div>
        
        <div class="section">
            <span class="label">%s结果:</span>
            <div class="content %s">%s</div>
        </div>
        
        <div style="text-align: center; margin-top: 30px;">
            <a href="/" class="button">返回首页</a>
            <button onclick="copyResult()" class="button button-secondary">复制结果</button>
        </div>
    </div>
    
    <script>
        function copyResult() {
            const resultText = document.querySelector('.result').textContent;
            navigator.clipboard.writeText(resultText).then(function() {
                alert('结果已复制到剪贴板');
            });
        }
    </script>
</body>
</html>`

	var resultClass, displayResult string
	if err != nil {
		resultClass = "error"
		displayResult = fmt.Sprintf("错误: %v", err)
	} else {
		resultClass = "result"
		displayResult = result
	}

	methodNames := map[string]string{
		"base64":           "Base64",
		"html_force":       "HTML实体编码(强制)",
		"html_force_hex":   "HTML实体编码(强制十六进制)",
		"html_special":     "HTML实体编码(特殊字符)",
		"url_force":        "URL编码(强制)",
		"url_special":      "URL编码(特殊字符)",
		"url_path_special": "URL路径编码(特殊字符)",
		"double_url":       "双重URL编码",
		"hex":              "十六进制编码",
		"unicode":          "Unicode中文编码",
		"md5":              "MD5",
		"sm3":              "SM3",
		"sha1":             "SHA1",
		"sha256":           "SHA-256",
		"sha512":           "SHA-512",
		"html":             "HTML",
		"url":              "URL",
		"url_path":         "URL路径",
	}

	methodName := methodNames[method]

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, html, opType, methodName, template.HTMLEscapeString(input), opType, resultClass, template.HTMLEscapeString(displayResult))
}

// 编码函数实现
func (tool *CodecTool) encodeBase64(input string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(input)), nil
}

func (tool *CodecTool) encodeHTMLForce(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		result.WriteString("&#" + strconv.Itoa(int(r)) + ";")
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeHTMLForceHex(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		result.WriteString("&#x" + strconv.FormatInt(int64(r), 16) + ";")
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeHTMLSpecial(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		switch r {
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '&':
			result.WriteString("&amp;")
		case '"':
			result.WriteString("&quot;")
		case '\'':
			result.WriteString("&#39;")
		default:
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeURLForce(input string) (string, error) {
	// 真正的强制URL编码 - 将所有字符都编码为URL格式
	var result strings.Builder
	for _, r := range input {
		// 将每个字符转换为UTF-8字节序列，然后对每个字节进行编码
		utf8Bytes := []byte(string(r))
		for _, b := range utf8Bytes {
			result.WriteString("%" + fmt.Sprintf("%02x", b))
		}
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeURLSpecial(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		if r <= 127 && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && !strings.ContainsRune("-._~", r) {
			result.WriteString("%" + fmt.Sprintf("%02X", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeURLPathSpecial(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		if r <= 127 && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && !strings.ContainsRune("-._~/", r) {
			result.WriteString("%" + fmt.Sprintf("%02X", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeDoubleURL(input string) (string, error) {
	first := url.QueryEscape(input)
	second := url.QueryEscape(first)
	return second, nil
}

func (tool *CodecTool) encodeHex(input string) (string, error) {
	return hex.EncodeToString([]byte(input)), nil
}

func (tool *CodecTool) encodeUnicode(input string) (string, error) {
	var result strings.Builder
	for _, r := range input {
		if r > 127 {
			result.WriteString("\\u" + fmt.Sprintf("%04X", r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func (tool *CodecTool) encodeMD5(input string) (string, error) {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func (tool *CodecTool) encodeSM3(input string) (string, error) {
	// 简单的哈希实现，替代SM3
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func (tool *CodecTool) encodeSHA1(input string) (string, error) {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func (tool *CodecTool) encodeSHA256(input string) (string, error) {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func (tool *CodecTool) encodeSHA512(input string) (string, error) {
	hash := sha512.Sum512([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

// 解码函数实现
func (tool *CodecTool) decodeBase64(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (tool *CodecTool) decodeHTML(input string) (string, error) {
	// 简单的HTML实体解码
	result := input

	// 解码数字实体
	result = tool.decodeHTMLEntities(result)

	return result, nil
}

func (tool *CodecTool) decodeHTMLEntities(input string) string {
	var result strings.Builder
	i := 0
	for i < len(input) {
		if input[i] == '&' {
			// 查找实体结束符
			end := strings.IndexByte(input[i:], ';')
			if end == -1 {
				result.WriteByte(input[i])
				i++
				continue
			}

			entity := input[i : i+end+1]
			decoded := tool.decodeSingleHTMLEntity(entity)
			if decoded != "" {
				result.WriteString(decoded)
				i += end + 1
			} else {
				result.WriteByte(input[i])
				i++
			}
		} else {
			result.WriteByte(input[i])
			i++
		}
	}
	return result.String()
}

func (tool *CodecTool) decodeSingleHTMLEntity(entity string) string {
	entity = entity[1 : len(entity)-1] // 移除 & 和 ;

	if strings.HasPrefix(entity, "#x") || strings.HasPrefix(entity, "#X") {
		// 十六进制
		hexStr := entity[2:]
		if val, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
			if r := rune(val); utf8.ValidRune(r) {
				return string(r)
			}
		}
	} else if strings.HasPrefix(entity, "#") {
		// 十进制
		decStr := entity[1:]
		if val, err := strconv.ParseInt(decStr, 10, 32); err == nil {
			if r := rune(val); utf8.ValidRune(r) {
				return string(r)
			}
		}
	} else {
		// 命名实体
		switch entity {
		case "lt":
			return "<"
		case "gt":
			return ">"
		case "amp":
			return "&"
		case "quot":
			return "\""
		case "apos":
			return "'"
		}
	}

	return ""
}

func (tool *CodecTool) decodeURL(input string) (string, error) {
	return url.QueryUnescape(input)
}

func (tool *CodecTool) decodeURLPath(input string) (string, error) {
	return url.PathUnescape(input)
}

func (tool *CodecTool) decodeDoubleURL(input string) (string, error) {
	// 先解码一次
	first, err := url.QueryUnescape(input)
	if err != nil {
		return "", err
	}
	// 再解码一次
	return url.QueryUnescape(first)
}

func (tool *CodecTool) decodeHex(input string) (string, error) {
	data, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (tool *CodecTool) decodeUnicode(input string) (string, error) {
	var result strings.Builder
	i := 0
	for i < len(input) {
		if input[i] == '\\' && i+1 < len(input) && (input[i+1] == 'u' || input[i+1] == 'U') {
			// 尝试解析Unicode转义序列
			hexStr := ""
			if input[i+1] == 'u' {
				if i+5 < len(input) {
					hexStr = input[i+2 : i+6]
					if len(hexStr) == 4 {
						if val, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
							if r := rune(val); utf8.ValidRune(r) {
								result.WriteRune(r)
								i += 6
								continue
							}
						}
					}
				}
			} else { // \U
				if i+9 < len(input) {
					hexStr = input[i+2 : i+10]
					if len(hexStr) == 8 {
						if val, err := strconv.ParseInt(hexStr, 16, 32); err == nil {
							if r := rune(val); utf8.ValidRune(r) {
								result.WriteRune(r)
								i += 10
								continue
							}
						}
					}
				}
			}
		}
		result.WriteByte(input[i])
		i++
	}
	return result.String(), nil
}
