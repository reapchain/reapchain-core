#!/bin/bash

cp -a ../rpc/openapi/ .vuepress/public/rpc/
git clone https://github.com/reapchain/spec.git specRepo && cp -r specRepo/spec . && rm -rf specRepo
