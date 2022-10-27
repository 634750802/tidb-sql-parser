import Highlight from 'highlight.js'
import sql from 'highlight.js/lib/languages/sql'
import js from 'highlight.js/lib/languages/javascript'
import darkTheme from "highlight.js/styles/dark.css?url";
import lightTheme from "highlight.js/styles/vs.css?url";

Highlight.registerLanguage('mysql', sql)
Highlight.registerLanguage('javascript', js)

if (typeof window !== 'undefined') {
    let link = document.createElement('link')
    link.href = darkTheme
    link.media = "(prefers-color-scheme: dark)"
    link.rel = 'stylesheet'
    document.head.append(link)

    link = document.createElement('link')
    link.href = lightTheme
    link.media = "(prefers-color-scheme: light)"
    link.rel = 'stylesheet'
    document.head.append(link)
}

export {default} from './Highlight.vue'
