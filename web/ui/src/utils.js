function formatHexDump(inputString) {
    if (inputString === undefined || inputString === null || inputString === '') {
        return ''
    }
// 创建一个 TextEncoder 对象
    const textEncoder = new TextEncoder();

// 将字符串编码为 Uint8Array
    const byteArray = textEncoder.encode(inputString);

    let result = '';
    for (let i = 0; i < byteArray.length; i += 16) {
        // 创建一行的十六进制字符串
        let hexLine = Array.from(byteArray.slice(i, i + 16), byte => byte.toString(16).padStart(2, '0')).join(' ');
        hexLine = hexLine.padEnd(48, ' ')

        // 创建一行的字符字符串
        const charLine = Array.from(byteArray.slice(i, i + 16), byte => {
            const char = String.fromCharCode(byte);
            // 只显示可打印字符，非可打印字符用点号表示
            return char.match(/[ -~]/) ? char : '.';
        }).join('');

        // 格式化行号
        const lineNumber = i.toString(16).padStart(8, '0');

        // 将一行的十六进制和字符字符串连接起来
        result += `${lineNumber}  ${hexLine}  ${charLine}\n`;
    }
    return result;
}


export {formatHexDump}


