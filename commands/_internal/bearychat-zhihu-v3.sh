#!/bin/bash

curl -s 'https://api.zhihu.com/search_v3?correction=1&excerpt_len=70&q=Bearychat&t=general' -H 'Authorization: Bearer 2.0AABA4akYAAAAMACRfmgrCgsAAABgAlVNsDYuWQDTS4Tg39E7mJAv_xlZXUvOCJ0BVA' --compressed \
| jq '[.data [1:] [] | {title: .highlight.title, author: .object.author.name, answer: .object.id, question: .object.question.id}]'