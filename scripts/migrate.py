#!/usr/bin/env python3

import argparse
import os
import re
import frontmatter

def get_md_files(content_dir: str) -> list[str]:
    md_files = []
    for root, _, files in os.walk(content_dir):
        for filename in files:
            if filename.endswith(".md"):
                md_files.append(os.path.join(root, filename))
    return md_files


def convert(filepath: str) -> frontmatter.Post:
    page = frontmatter.load(filepath)

    try:
        page["aliases"] = page["redirect_from"]
        del page["redirect_from"]
    except:
        pass

    # handle CLI reference stubs
    if "/engine/reference/commandline/" in filepath:
        try:
            if page["datafolder"]:
                page = convert_cli(page)
                return page
        except:
            pass

    # handle sample stubs
    if "/samples/" in filepath:
        page = convert_sample(page)
        return page

    # all other files
    convert_other(page)
    return page


def convert_cli(page: frontmatter.Post) -> frontmatter.Post:
    page["layout"] = "cli"
    page.content = re.sub(r"\{% include cli.*", "", page.content)
    page.content = re.sub(
        r"\{% include (.*(md|html)) %}", r'{{< include "\1" >}}', page.content
    )
    return page


def convert_sample(page: frontmatter.Post) -> frontmatter.Post:
    page.content = re.sub(r"\{% include_relative samples_body.*", "", page.content)
    return page


def convert_other(page: frontmatter.Post) -> frontmatter.Post:
    page.content = re.sub(r"\{:\s*target=(.|\n)*?}", "", page.content)
    page.content = re.sub(r"\{%-?\s*include eula(.|\n)*?}", "", page.content)
    page.content = re.sub(r"\{%\sraw(.|\n)*?}", "", page.content)
    page.content = re.sub(r"\{%\sendraw(.|\n)*?}", "", page.content)
    page.content = re.sub(
        r"\{:\s*\.(important|warning|tip|experimental|restricted).*?}",
        r"{ .\1 }",
        page.content,
    )
    page.content = re.sub(
        r"\{:\s*(\.invertible|\.text-center|style=|width=|height=|class=|\.accept-eula).*?}",
        "",
        page.content,
    )
    page.content = re.sub(
        r"\{% include (.*(md|html)) %}", r'{{< include "\1" >}}', page.content
    )
    page.content = re.sub(
        r"!\[(.*?)]\((.*?)\)\{:\s*.inline\s*}",
        r'{{< inline-image src="\2" alt="\1" >}}',
        page.content,
    )
    page.content = re.sub(
        r"\[(.*)?]\((.*?)\)\{:\s*\.button.*}",
        r'{{< button text="\1" url="\2" >}}',
        page.content,
    )
    page.content = re.sub(
        r"\{% include desktop-install\.md (.*) %}",
        r"{{< desktop-install \1 >}}",
        page.content,
    )
    page.content = re.sub(
        r"\{% include release-date\.html (.*) %}",
        r"{{< release-date \1 >}}",
        page.content,
    )
    page.content = re.sub(
        r"\{% include (admin-(.*)?)\.md (.*) %}",
        r"{{% admin-\2 \3 %}}",
        page.content,
    )
    page.content = re.sub(
        r"\{\{\s*site\.(.*?)\s*}}", r'{{% param "\1" %}}', page.content
    )
    page.content = re.sub(r'\{:\s*id="(.*)" ?}', r"{ #\1 }", page.content)

    # handle inline assign
    assign_pattern = re.compile(r"{% assign (.*) = (['\"].*) %}\n?")
    assign = assign_pattern.search(page.content)
    if assign:
        key = assign.group(1)
        value = assign.group(2)
        page[key] = value[1:-1]
        page.content = assign_pattern.sub("", page.content)
        page.content = re.sub(
            rf'{{{{\s*{key}\s*}}}}', rf'{{{{% param "{key}" %}}}}', page.content
        )

    return page


def write_results(filepath: str, page: frontmatter.Post) -> None:
    f = open(filepath, "w")
    if "/includes/" in filepath:
        # skip front matter for includes
        data = page.content
    else:
        data = frontmatter.dumps(page, sort_keys=False)
    f.write(data)
    f.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="Migrate",
        description="Convert content from Jekyll to Hugo",
    )
    parser.add_argument("-d", type=str, help="Content directory")
    parser.add_argument("-f", type=str, help="Single input file")
    args = parser.parse_args()

    if args.f:
        converted_page = convert(args.f)
        write_results(args.f, converted_page)
        os._exit(os.EX_OK)

    content_dir = os.path.join(os.getcwd(), "content")
    if args.d:
        content_dir = os.path.abspath(args.d)

    md_files = get_md_files(content_dir)

    for filepath in md_files:
        converted_page = convert(filepath)
        write_results(filepath, converted_page)

    os._exit(os.EX_OK)
