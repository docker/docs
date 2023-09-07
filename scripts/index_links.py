#!/usr/bin/env python3

import os
import re


def get_md_files(content_dir: str) -> list[str]:
    md_files = []
    for root, _, files in os.walk(content_dir):
        for filename in files:
            if filename.endswith(".md"):
                md_files.append(os.path.join(root, filename))
    return md_files


if __name__ == "__main__":
    content_dir = os.path.join(os.getcwd(), "content")
    md_files = get_md_files(content_dir)

    for filepath in md_files:
        file = open(filepath, "r")
        content = file.read()
        file.close()

        content = re.sub(r"([^-_])index\.md", r"\1_index.md", content, flags=re.M)
        
        file = open(filepath, "w")
        file.write(content)
        file.close()

    os._exit(os.EX_OK)
