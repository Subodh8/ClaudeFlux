import zipfile
import xml.etree.ElementTree as ET

def extract_text_from_docx(file_path):
    z = zipfile.ZipFile(file_path)
    xml_content = z.read('word/document.xml')
    root = ET.fromstring(xml_content)
    ns = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
    text = []
    for node in root.iter(f'{{{ns["w"]}}}t'):
        if node.text:
            text.append(node.text)
    return '\n'.join(text)

if __name__ == '__main__':
    text = extract_text_from_docx(r'c:\Users\subod\Desktop\Github\ClaudeFlux\ClaudeFlux_Repository_Blueprint.docx')
    with open(r'c:\Users\subod\Desktop\Github\ClaudeFlux\docx_content.txt', 'w', encoding='utf-8') as f:
        f.write(text)
