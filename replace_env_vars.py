import os
import re

directory = '/home/ubuntu/develop/birds_eye_v3'
directories_to_ignore = {'.git', 'node_modules', 'vendor', 'dist', 'build'}

replacements = {
    'MYSQL_ROOT_PASSWORD': 'BIRDSEYE_MYSQL_ROOT_PASSWORD',
    'GO_API_PORT': 'BIRDSEYE_GO_API_PORT',
    'BIRDSEYEAPI_EXECUTION_MODE': 'BIRDSEYE_BIRDSEYEAPI_EXECUTION_MODE',
    'OPENAI_MODEL': 'BIRDSEYE_OPENAI_MODEL',
    'OPENAI_CHAT_ENDPOINT': 'BIRDSEYE_OPENAI_CHAT_ENDPOINT',
    'BIRDSEYEAPI_V2_OPENAI_API_KEY': 'BIRDSEYE_BIRDSEYEAPI_V2_OPENAI_API_KEY',
    'BIRDSEYEAPI_V2_CLAUDE_API_KEY': 'BIRDSEYE_BIRDSEYEAPI_V2_CLAUDE_API_KEY',
    'SCRAPING_ARTICLES': 'BIRDSEYE_SCRAPING_ARTICLES',
    'SELENIUM_URL': 'BIRDSEYE_SELENIUM_URL'
}

for root, dirs, files in os.walk(directory):
    dirs[:] = [d for d in dirs if d not in directories_to_ignore]
    for file in files:
        filepath = os.path.join(root, file)
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
        except (UnicodeDecodeError, IsADirectoryError):
            continue
            
        new_content = content
        for k, v in replacements.items():
            new_content = re.sub(r'(?<!BIRDSEYE_)(?<!BIRDSEYE_V2_)\b' + re.escape(k) + r'\b', v, new_content)
            
        if new_content != content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(new_content)
            print(f"Updated {filepath}")
