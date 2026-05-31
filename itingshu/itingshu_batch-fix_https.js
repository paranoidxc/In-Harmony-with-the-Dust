// ==UserScript==
// @name         听书网批量下载 - 可选起始集数版
// @namespace    http://tampermonkey.net/
// @version      8.0
// @description  智能识别列表，支持选择起始集数，连续跳转下载，并兼容 mp3/m4a 音频
// @match        https://www.itingshu.net/itingshus/*
// @match        https://www.itingshu.net/play/*
// @grant        GM_download
// @run-at       document-start
// ==/UserScript==

(function() {
    'use strict';

    const WAIT_SECONDS = 15;   // 可在此修改默认等待秒数
    const MAX_WAIT = 10000;
    const CHECK_INTERVAL = 500;
    const PAGE_API_WAIT_MS = 2500;
    const DOWNLOAD_SETTLE_TIMEOUT_MS = 30000;
    const PAGE_API_RESULT_EVENT = 'itingshu:page-api-result';
    const PAGE_API_MESSAGE_TYPE = 'itingshu:page-api-result';
    const processedApiEventIds = new Set();
    const AUDIO_URL_RE = /\.((mp3)|(m4a))(\?|#|$)/i;
    let pageApiAudioResolved = false;
    let audioDetected = false;
    let pageApiWaiters = [];
    let currentWaitSeconds = WAIT_SECONDS;
    let activeDownload = null;

    function resetPlayPageState() {
        pageApiAudioResolved = false;
        audioDetected = false;
        pageApiWaiters = [];
        processedApiEventIds.clear();
        activeDownload = null;
    }

    function createActiveDownload(url, fileName) {
        let settle;
        const settled = new Promise((resolve) => {
            settle = resolve;
        });

        activeDownload = {
            url: url,
            fileName: fileName,
            startedAt: Date.now(),
            finished: false,
            settle: (result) => {
                if (activeDownload && activeDownload.url === url && !activeDownload.finished) {
                    activeDownload.finished = true;
                    settle(result);
                }
            },
            settled: settled
        };

        return activeDownload;
    }

    function waitForDomReady() {
        if (document.readyState === 'loading') {
            return new Promise((resolve) => {
                document.addEventListener('DOMContentLoaded', resolve, { once: true });
            });
        }
        return Promise.resolve();
    }

    function isAudioUrl(url) {
        return typeof url === 'string' && AUDIO_URL_RE.test(url);
    }

    function getAudioExtension(url, fallback = 'mp3') {
        const matched = typeof url === 'string' ? url.match(/\.((mp3)|(m4a))(\?|#|$)/i) : null;
        return matched ? matched[1].toLowerCase() : fallback;
    }

    function injectPageRequestDebugHooks() {
        const script = document.createElement('script');
        script.textContent = `(() => {
            if (window.__itingshuDebugHookInstalled) return;
            window.__itingshuDebugHookInstalled = true;

            const previewText = (text, limit = 500) => {
                if (!text) return '';
                return text.length > limit ? text.slice(0, limit) + '...' : text;
            };
            const isTarget = (url) => typeof url === 'string' && url.includes('/api/mapi/play');
            const extractAudioUrl = (payload) => {
                const visited = new Set();

                const walk = (value) => {
                    if (!value) return null;
                    if (typeof value === 'string') {
                        return /\\.((mp3)|(m4a))(\\?|#|$)/i.test(value) ? value : null;
                    }
                    if (typeof value !== 'object') return null;
                    if (visited.has(value)) return null;
                    visited.add(value);

                    if (typeof value.url === 'string' && /\\.((mp3)|(m4a))(\\?|#|$)/i.test(value.url)) return value.url;
                    if (typeof value.src === 'string' && /\\.((mp3)|(m4a))(\\?|#|$)/i.test(value.src)) return value.src;
                    if (typeof value.file === 'string' && /\\.((mp3)|(m4a))(\\?|#|$)/i.test(value.file)) return value.file;

                    const values = Array.isArray(value) ? value : Object.values(value);
                    for (const item of values) {
                        const found = walk(item);
                        if (found) return found;
                    }
                    return null;
                };

                return walk(payload);
            };
            const parseResponse = (responseText) => {
                if (!responseText) return { data: null, audioUrl: null };
                try {
                    const data = JSON.parse(responseText);
                    return { data, audioUrl: extractAudioUrl(data) };
                } catch (error) {
                    const matchedUrl = responseText.match(/https?:\\/\\/[^"'\\\\\\s]+\\.((mp3)|(m4a))[^"'\\\\\\s]*/i);
                    return { data: null, audioUrl: matchedUrl ? matchedUrl[0] : null };
                }
            };
            const notifyAudio = (meta) => {
                const payload = {
                    ...meta,
                    eventId: meta.eventId || String(Date.now()) + '-' + Math.random().toString(36).slice(2)
                };
                window.dispatchEvent(new CustomEvent('${PAGE_API_RESULT_EVENT}', {
                    detail: payload
                }));
                window.postMessage({
                    source: 'itingshu-userscript-hook',
                    type: '${PAGE_API_MESSAGE_TYPE}',
                    payload: payload
                }, '*');
            };

            if (typeof window.fetch === 'function') {
                const originalFetch = window.fetch;
                window.fetch = async function(...args) {
                    const input = args[0];
                    const init = args[1] || {};
                    const url = typeof input === 'string' ? input : (input && input.url) || '';
                    const method = (init.method || (input && input.method) || 'GET').toUpperCase();
                    const body = init.body;
                    const response = await originalFetch.apply(this, args);

                    if (isTarget(url)) {
                        try {
                            const cloned = response.clone();
                            const text = await cloned.text();
                            const parsed = parseResponse(text);
                            console.log('[页面请求] FETCH', {
                                method,
                                url,
                                body: typeof body === 'string' ? body : '',
                                status: response.status,
                                response: previewText(text)
                            });
                            if (parsed.audioUrl) {
                                notifyAudio({
                                    transport: 'fetch',
                                    method,
                                    url,
                                    body: typeof body === 'string' ? body : '',
                                    title: parsed.data && parsed.data.name ? parsed.data.name : '',
                                    audioUrl: parsed.audioUrl
                                });
                            }
                        } catch (error) {
                            console.warn('[页面请求] FETCH 响应读取失败', { method, url, error });
                        }
                    }

                    return response;
                };
            }

            const originalOpen = XMLHttpRequest.prototype.open;
            const originalSend = XMLHttpRequest.prototype.send;

            XMLHttpRequest.prototype.open = function(method, url, ...rest) {
                this.__itingshuMethod = String(method || 'GET').toUpperCase();
                this.__itingshuUrl = String(url || '');
                return originalOpen.call(this, method, url, ...rest);
            };

            XMLHttpRequest.prototype.send = function(body) {
                if (isTarget(this.__itingshuUrl)) {
                    this.addEventListener('loadend', function() {
                        try {
                            const responseText = this.responseText || '';
                            const parsed = parseResponse(responseText);
                            console.log('[页面请求] XHR', {
                                method: this.__itingshuMethod,
                                url: this.__itingshuUrl,
                                body: typeof body === 'string' ? body : '',
                                status: this.status,
                                response: previewText(responseText)
                            });
                            if (parsed.audioUrl) {
                                notifyAudio({
                                    transport: 'xhr',
                                    method: this.__itingshuMethod,
                                    url: this.__itingshuUrl,
                                    body: typeof body === 'string' ? body : '',
                                    title: parsed.data && parsed.data.name ? parsed.data.name : '',
                                    audioUrl: parsed.audioUrl
                                });
                            }
                        } catch (error) {
                            console.warn('[页面请求] XHR 响应读取失败', {
                                method: this.__itingshuMethod,
                                url: this.__itingshuUrl,
                                error
                            });
                        }
                    }, { once: true });
                }

                return originalSend.call(this, body);
            };
        })();`;

        const mount = document.documentElement || document.head;
        if (mount) {
            mount.appendChild(script);
            script.remove();
        } else {
            document.addEventListener('readystatechange', () => {
                const target = document.documentElement || document.head;
                if (!target || script.isConnected) return;
                target.appendChild(script);
                script.remove();
            }, { once: true });
        }
    }

    injectPageRequestDebugHooks();

    function listenPageApiResults() {
        const handleDetail = (detail) => {
            if (!detail.audioUrl) return;
            if (detail.eventId && processedApiEventIds.has(detail.eventId)) return;
            if (detail.eventId) processedApiEventIds.add(detail.eventId);
            pageApiAudioResolved = true;
            audioDetected = true;
            for (const resolve of pageApiWaiters) resolve(true);
            pageApiWaiters = [];
            const task = loadTask();
            const currentTitle = task?.links?.[task.currentIndex]?.title || document.title;
            const title = detail.title || currentTitle;
            console.log(`[页面接口复用] ${detail.transport?.toUpperCase() || 'UNKNOWN'} ${detail.method || ''} 命中音频地址`);
            downloadMP3(detail.audioUrl, title);
        };

        window.addEventListener(PAGE_API_RESULT_EVENT, (event) => {
            handleDetail(event.detail || {});
        });

        window.addEventListener('message', (event) => {
            if (event.source !== window) return;
            const data = event.data || {};
            if (data.source !== 'itingshu-userscript-hook' || data.type !== PAGE_API_MESSAGE_TYPE) return;
            handleDetail(data.payload || {});
        });
    }

    listenPageApiResults();

    function waitForPageApiAudio(timeoutMs = PAGE_API_WAIT_MS) {
        if (pageApiAudioResolved || audioDetected) {
            return Promise.resolve(true);
        }

        return new Promise((resolve) => {
            const done = (matched) => {
                clearTimeout(timer);
                pageApiWaiters = pageApiWaiters.filter((item) => item !== done);
                resolve(Boolean(matched));
            };

            const timer = setTimeout(() => {
                done(pageApiAudioResolved || audioDetected);
            }, timeoutMs);

            pageApiWaiters.push(done);
        });
    }

    // ========== 智能获取播放列表链接 ==========
    async function getPlaylistLinksWithRetry() {
        const startTime = Date.now();
        let container = null;
        while (Date.now() - startTime < MAX_WAIT) {
            container = document.querySelector('#playlist, .playlist, div[class*="playlist"], ul[class*="playlist"]');
            if (container) break;
            await new Promise(r => setTimeout(r, CHECK_INTERVAL));
        }

        function extractLinks(root) {
            const items = [];
            const links = root.querySelectorAll('a[href*="/play/"]');
            for (const link of links) {
                const url = link.href;
                if (url && !items.some(l => l.url === url)) {
                    let title = link.textContent.trim();
                    title = title.replace(/^\d{4}-\d{2}-\d{2}\s*/, '');
                    title = title.replace(/\s*(在线播放|免费收听)$/, '');
                    if (title) items.push({ url, title });
                }
            }
            return items;
        }

        let items = [];
        if (container) {
            items = extractLinks(container);
            console.log(`[抓取] 从容器中找到 ${items.length} 个音频`);
        }
        if (items.length === 0) {
            items = extractLinks(document);
            console.log(`[抓取] 全页搜索找到 ${items.length} 个音频`);
        }
        return items;
    }

    // ========== 手动选择容器（备用） ==========
    function manualSelectContainer() {
        return new Promise((resolve) => {
            const overlay = document.createElement('div');
            overlay.style.cssText = 'position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.5);z-index:100000;cursor:crosshair;';
            const tip = document.createElement('div');
            tip.textContent = '请点击包含播放列表的区域（例如“051_蛋挞女士_上”所在的容器）';
            tip.style.cssText = 'position:fixed;top:20px;left:50%;transform:translateX(-50%);background:#fff;color:#000;padding:10px 20px;border-radius:8px;z-index:100001;';
            document.body.appendChild(overlay);
            document.body.appendChild(tip);

            overlay.addEventListener('click', (e) => {
                let target = e.target;
                while (target && target !== document.body) {
                    const links = target.querySelectorAll('a[href*="/play/"]');
                    if (links.length > 0) {
                        const items = [];
                        for (const link of links) {
                            const url = link.href;
                            if (url && !items.some(l => l.url === url)) {
                                let title = link.textContent.trim().replace(/^\d{4}-\d{2}-\d{2}\s*/, '');
                                if (title) items.push({ url, title });
                            }
                        }
                        if (items.length > 0) {
                            overlay.remove();
                            tip.remove();
                            resolve(items);
                            return;
                        }
                    }
                    target = target.parentElement;
                }
                alert('未找到有效列表，请重新点击包含播放链接的区域');
            });
        });
    }

    // ========== 状态管理（连续跳转用） ==========
    function saveTask(links, startIndex, listPageUrl) {
        const task = {
            links: links.map(l => ({ url: l.url, title: l.title })),
            currentIndex: startIndex,
            total: links.length,
            listPageUrl: listPageUrl
        };
        sessionStorage.setItem('dl_task', JSON.stringify(task));
    }
    function loadTask() {
        const raw = sessionStorage.getItem('dl_task');
        return raw ? JSON.parse(raw) : null;
    }
    function clearTask() {
        sessionStorage.removeItem('dl_task');
    }

    // ========== 下载与监听 ==========
    function downloadMP3(url, title) {
        if (!isAudioUrl(url)) return;
        // 将下载地址中的 http 替换为 https
        url = url.replace(/^http:\/\//i, 'https://');
        audioDetected = true;
        const extension = getAudioExtension(url);
        let fileName = title.replace(/[\\/:*?"<>|]/g, '_');
        // 将文件名中的所有下划线替换为短横线
        fileName = fileName.replace(/_/g, '-');
        fileName = fileName.replace(/\.(mp3|m4a)$/i, '');
        fileName += `.${extension}`;

        if (activeDownload && !activeDownload.finished) {
            if (activeDownload.url === url || activeDownload.fileName === fileName) {
                console.log(`[下载] 忽略同页重复触发: ${fileName}`);
                return activeDownload.settled;
            }
            console.warn(`[下载] 当前页已有进行中的下载，忽略额外触发: ${fileName}`);
            return activeDownload.settled;
        }

        const downloadTask = createActiveDownload(url, fileName);
        console.log(`[下载] ${fileName}`);
        GM_download({
            url: url,
            name: fileName,
            saveAs: false,
            conflictAction: 'uniquify',
            onload: () => {
                console.log(`✅ 成功: ${fileName}`);
                downloadTask.settle({ ok: true, fileName: fileName });
            },
            onerror: (err) => {
                console.error(`❌ 失败: ${fileName}`, err);
                downloadTask.settle({ ok: false, fileName: fileName, error: err });
                window.open(url, '_blank');
            }
        });

        return downloadTask.settled;
    }

    let audioMonitor = null;
    function startAudioMonitor(currentTitle) {
        if (audioMonitor) audioMonitor.disconnect();
        const captured = new Set();
        audioMonitor = new PerformanceObserver((list) => {
            for (const entry of list.getEntries()) {
                if (isAudioUrl(entry.name) && !captured.has(entry.name)) {
                    captured.add(entry.name);
                    downloadMP3(entry.name, currentTitle);
                }
            }
        });
        audioMonitor.observe({ entryTypes: ['resource'] });
    }

    function previewText(text, limit = 300) {
        if (!text) return '';
        return text.length > limit ? `${text.slice(0, limit)}...` : text;
    }

    function extractAudioUrl(payload) {
        const visited = new Set();

        function walk(value) {
            if (!value) return null;
            if (typeof value === 'string') {
                return isAudioUrl(value) ? value : null;
            }
            if (typeof value !== 'object') return null;
            if (visited.has(value)) return null;
            visited.add(value);

            if (typeof value.url === 'string' && isAudioUrl(value.url)) return value.url;
            if (typeof value.src === 'string' && isAudioUrl(value.src)) return value.src;
            if (typeof value.file === 'string' && isAudioUrl(value.file)) return value.file;

            const values = Array.isArray(value) ? value : Object.values(value);
            for (const item of values) {
                const found = walk(item);
                if (found) return found;
            }
            return null;
        }

        return walk(payload);
    }

    function parseAudioResponse(responseText) {
        if (!responseText) return { raw: responseText, data: null, audioUrl: null };

        try {
            const data = JSON.parse(responseText);
            return {
                raw: responseText,
                data: data,
                audioUrl: extractAudioUrl(data)
            };
        } catch (error) {
            const matchedUrl = responseText.match(/https?:\/\/[^"'\\\s]+\.((mp3)|(m4a))[^"'\\\s]*/i);
            return {
                raw: responseText,
                data: null,
                audioUrl: matchedUrl ? matchedUrl[0] : null
            };
        }
    }

    function requestMetaApiWithXhr(url, method, body) {
        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();
            xhr.open(method, url, true);
            xhr.withCredentials = true;
            xhr.responseType = 'text';
            xhr.setRequestHeader('Accept', 'application/json, text/plain, */*');
            xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
            if (method === 'POST') {
                xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
            }

            xhr.onload = () => {
                resolve({
                    ok: xhr.status >= 200 && xhr.status < 300,
                    status: xhr.status,
                    text: xhr.responseText || ''
                });
            };
            xhr.onerror = () => reject(new Error(`XHR network error (${method})`));
            xhr.ontimeout = () => reject(new Error(`XHR timeout (${method})`));
            xhr.send(method === 'POST' ? body : null);
        });
    }

    async function requestMetaApiWithFetch(url, method, body) {
        const headers = {
            'Accept': 'application/json, text/plain, */*',
            'X-Requested-With': 'XMLHttpRequest'
        };
        if (method === 'POST') {
            headers['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
        }

        const response = await fetch(url, {
            method: method,
            headers: headers,
            body: method === 'POST' ? body : undefined,
            credentials: 'same-origin'
        });

        return {
            ok: response.ok,
            status: response.status,
            text: await response.text()
        };
    }

    async function tryDownloadFromMetaApi(currentTitle) {
        if (pageApiAudioResolved || audioDetected) {
            console.log('[兼容下载] 已通过页面接口或资源监听拿到音频，跳过补发请求');
            return true;
        }

        const nid = document.querySelector('meta[name="_b"]')?.content?.trim();
        const cid = document.querySelector('meta[name="_p"]')?.content?.trim();
        if (!nid || !cid) {
            console.log('[兼容下载] 未找到 _b / _p 元数据');
            return false;
        }

        const params = new URLSearchParams({
            nid: nid,
            cid: cid,
            sort: 'read'
        });

        console.log(`[兼容下载] 使用 nid=${nid}, cid=${cid}`);

        const requestOptions = [
            {
                transport: 'xhr',
                method: 'POST'
            },
            {
                transport: 'fetch',
                method: 'POST'
            },
            {
                transport: 'xhr',
                method: 'GET'
            },
            {
                transport: 'fetch',
                method: 'GET'
            }
        ];

        for (const options of requestOptions) {
            const url = options.method === 'GET'
                ? `https://www.itingshu.net/api/mapi/play?${params.toString()}`
                : 'https://www.itingshu.net/api/mapi/play';
            const body = params.toString();

            try {
                const response = options.transport === 'xhr'
                    ? await requestMetaApiWithXhr(url, options.method, body)
                    : await requestMetaApiWithFetch(url, options.method, body);

                if (!response.ok) {
                    console.warn(`[兼容下载] ${options.transport.toUpperCase()} ${options.method} 请求失败: ${response.status}`, previewText(response.text));
                    continue;
                }

                const parsed = parseAudioResponse(response.text);
                if (parsed.audioUrl) {
                    const preferredTitle = parsed.data && parsed.data.name ? parsed.data.name : currentTitle;
                    console.log(`[兼容下载] ${options.transport.toUpperCase()} ${options.method} 命中音频地址`);
                    downloadMP3(parsed.audioUrl, preferredTitle || document.title);
                    return true;
                }

                console.warn(
                    `[兼容下载] ${options.transport.toUpperCase()} ${options.method} 返回中未找到音频地址`,
                    parsed.data || previewText(parsed.raw)
                );
            } catch (error) {
                console.warn(`[兼容下载] ${options.transport.toUpperCase()} ${options.method} 请求异常`, error);
            }
        }

        return false;
    }

    async function main() {
        // ========== 播放页逻辑（连续跳转） ==========
        if (window.location.href.includes('/play/')) {
            resetPlayPageState();
            const task = loadTask();
            if (!task) {
                console.error('没有任务信息，请从列表页开始');
            } else {
                const currentIdx = task.currentIndex;
                const links = task.links;
                if (currentIdx >= links.length) {
                    clearTask();
                    alert(`✅ 共 ${links.length} 集全部下载完成！`);
                    window.location.href = task.listPageUrl;
                } else {
                    const current = links[currentIdx];
                    console.log(`[播放页] 第 ${currentIdx+1}/${links.length}: ${current.title}`);
                    startAudioMonitor(current.title);
                    void (async () => {
                        const matched = await waitForPageApiAudio();
                        if (matched) {
                            console.log('[播放页] 已通过页面接口或资源监听拿到音频，跳过兼容请求');
                            return;
                        }
                        await tryDownloadFromMetaApi(current.title);
                    })();

                    let seconds = currentWaitSeconds;
                    const timerDiv = document.createElement('div');
                    timerDiv.style.cssText = 'position:fixed; top:20px; right:20px; z-index:99999; background:#ff9800; padding:12px 20px; border-radius:8px; font-weight:bold; font-size:14px;';
                    timerDiv.textContent = `📖 ${current.title} (${currentIdx+1}/${links.length}) | ⏱️ ${seconds}秒后下一集`;
                    document.body.appendChild(timerDiv);

                    async function goNext() {
                        if (activeDownload && !activeDownload.finished) {
                            timerDiv.textContent = `📖 ${current.title} (${currentIdx+1}/${links.length}) | ⏳ 等待下载完成`;
                            const result = await Promise.race([
                                activeDownload.settled,
                                new Promise((resolve) => {
                                    setTimeout(() => resolve({
                                        ok: false,
                                        timeout: true,
                                        fileName: activeDownload ? activeDownload.fileName : ''
                                    }), DOWNLOAD_SETTLE_TIMEOUT_MS);
                                })
                            ]);

                            if (result && result.timeout) {
                                console.warn(`[播放页] 下载等待超时，继续下一集: ${result.fileName || current.title}`);
                            }
                        }

                        timerDiv.remove();
                        const nextIdx = currentIdx + 1;
                        if (nextIdx < links.length) {
                            task.currentIndex = nextIdx;
                            saveTask(links, nextIdx, task.listPageUrl);
                            window.location.href = links[nextIdx].url;
                        } else {
                            clearTask();
                            alert(`✅ 共 ${links.length} 集全部下载完成！`);
                            window.location.href = task.listPageUrl;
                        }
                    }

                    const timer = setInterval(() => {
                        seconds--;
                        if (timerDiv) timerDiv.textContent = `📖 ${current.title} (${currentIdx+1}/${links.length}) | ⏱️ ${seconds}秒后下一集`;
                        if (seconds <= 0) {
                            clearInterval(timer);
                            void goNext();
                        }
                    }, 1000);
                }
            }
        }

        // ========== 列表页逻辑（带起始集数选择） ==========
        else if (window.location.href.includes('/itingshus/')) {
            let links = await getPlaylistLinksWithRetry();
            if (links.length === 0) {
                // 手动选择
                const errorDiv = document.createElement('div');
                errorDiv.innerHTML = `
                    <div style="position:fixed; bottom:20px; right:20px; z-index:99999; background:#f44336; color:white; padding:15px; border-radius:10px;">
                        <b>⚠️ 未自动识别到列表</b><br>
                        <button id="auto-retry" style="margin-top:8px;">🔄 重试</button>
                        <button id="manual-select" style="margin-top:8px; margin-left:8px;">✋ 手动指定</button>
                    </div>
                `;
                document.body.appendChild(errorDiv);
                document.getElementById('auto-retry').onclick = () => location.reload();
                document.getElementById('manual-select').onclick = async () => {
                    errorDiv.remove();
                    links = await manualSelectContainer();
                    if (links.length) showStartDialog(links);
                    else alert('识别失败，请刷新页面重试');
                };
                return;
            }
            showStartDialog(links);
        }

        function showStartDialog(links) {
            // 生成列表预览（前10集）
            let preview = '<div style="max-height:150px; overflow-y:auto; margin-top:8px; font-size:12px;">';
            links.slice(0, 10).forEach((l, idx) => {
                preview += `<div>${idx+1}. ${l.title.substring(0, 35)}</div>`;
            });
            if (links.length > 10) preview += `<div>... 共${links.length}集</div>`;
            preview += '</div>';

            // 创建对话框（使用 div 模拟）
            const dialog = document.createElement('div');
            dialog.id = 'start-dialog';
            dialog.innerHTML = `
                <div style="position:fixed; top:80%; left:50%; transform:translate(-50%,-50%); z-index:100000; background:#1e1e1e; color:white; padding:20px; border-radius:12px; width:350px; max-width:90%; box-shadow:0 4px 20px black; font-family:monospace;">
                    <h3 style="margin:0 0 10px;">🎵 批量下载</h3>
                    <div>找到 <span style="color:#4caf50;">${links.length}</span> 个音频</div>
                    <div style="margin:10px 0;">
                        <label>▶️ 从第 <input type="number" id="startIndexInput" min="1" max="${links.length}" value="1" style="width:60px;"> 集开始</label>
                        <span style="font-size:12px; margin-left:8px;">(1-${links.length})</span>
                    </div>
                    <div style="margin:10px 0;">
                        ⏱️ 每集等待 <input type="number" id="waitSecondsInput" min="5" max="60" value="${currentWaitSeconds}" style="width:50px;"> 秒
                    </div>
                    ${preview}
                    <div style="display:flex; gap:10px; margin-top:15px;">
                        <button id="confirmStart" style="flex:1; padding:8px; background:#4caf50;">✅ 开始下载</button>
                        <button id="cancelStart" style="flex:1; padding:8px; background:#555;">❌ 取消</button>
                    </div>
                </div>
            `;
            document.body.appendChild(dialog);

            document.getElementById('confirmStart').onclick = () => {
                let start = parseInt(document.getElementById('startIndexInput').value);
                if (isNaN(start)) start = 1;
                start = Math.min(Math.max(start, 1), links.length);
                const waitSec = parseInt(document.getElementById('waitSecondsInput').value);
                if (!isNaN(waitSec) && waitSec > 0) currentWaitSeconds = waitSec;

                // 检查是否有未完成任务（可选）
                const existing = loadTask();
                if (existing && existing.links.length === links.length && existing.links[0].url === links[0].url) {
                    if (confirm(`检测到未完成的任务（上次进度：${existing.currentIndex+1}/${existing.total}），是否继续上次？\n点击“确定”继续，取消则从新选择的第 ${start} 集开始。`)) {
                        window.location.href = existing.links[existing.currentIndex].url;
                        dialog.remove();
                        return;
                    } else {
                        clearTask();
                    }
                }

                const startIndex = start - 1;
                saveTask(links, startIndex, window.location.href);
                dialog.remove();
                window.location.href = links[startIndex].url;
            };
            document.getElementById('cancelStart').onclick = () => {
                dialog.remove();
            };
        }
    }

    void waitForDomReady().then(main);
})();
