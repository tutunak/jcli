// Minimal shell syntax highlighter.
//
// Deliberately not a library: the page needs one language, in a handful of blocks,
// and a dependency would outweigh the problem. Blocks opt in with data-lang="shell"
// and hold plain text in the HTML, so crawlers, language models and no-JS visitors
// read the code unchanged -- colour is the only thing this adds.
//
// One ordered alternation does the whole job. Order is the correctness argument:
// comments and strings match first, so nothing inside them can be re-tokenized.
// The lookbehinds need a modern engine. If one refuses the pattern, leave the code
// as the plain text it already is in the HTML rather than marking it up halfway.
function buildShellRe(parts) {
    try {
        return new RegExp(parts.map(r => r.source).join('|'), 'g');
    } catch {
        return null;
    }
}

const SHELL_RE = buildShellRe([
    // # to end of line
    /(?<comment>#[^\n]*)/,
    // '...' takes no escapes in POSIX shell; "..." does. Neither may cross a
    // newline, so an unbalanced quote discolours one line instead of the rest
    // of the block.
    /(?<string>'[^'\n]*'|"(?:[^"\\\n]|\\.)*")/,
    // name immediately followed by () -- a function definition
    /(?<fn>\b[A-Za-z_][\w-]*(?=\s*\(\s*\)))/,
    // ${...}, $name, and the usual specials
    /(?<variable>\$\{[^}]*\}|\$[A-Za-z_]\w*|\$[@*#?$!0-9-])/,
    // The lookbehind keeps a word from lighting up inside a path, URL or module
    // spec: /usr/local/bin must not read "local" as a keyword, and
    // github.com/tutunak/jcli@latest must not read "jcli" as a command.
    /(?<keyword>(?<![\w./@-])(?:if|then|else|elif|fi|for|while|until|do|done|case|esac|in|function|return|local|export|readonly|shift|alias|unalias|source)\b)/,
    /(?<builtin>(?<![\w./@-])(?:jcli|git|printf|echo|read|cd|head|tail|grep|sed|awk|cut|sort|wc|tr|xargs)\b)/,
]);

const ESCAPES = { '&': '&amp;', '<': '&lt;', '>': '&gt;' };

function escapeHtml(s) {
    return s.replace(/[&<>]/g, c => ESCAPES[c]);
}

function highlightShell(code) {
    if (!SHELL_RE) return escapeHtml(code);

    let out = '';
    let last = 0;
    let m;

    SHELL_RE.lastIndex = 0;
    while ((m = SHELL_RE.exec(code)) !== null) {
        // Every alternative requires at least one character, so lastIndex always
        // advances and this cannot spin.
        const entry = Object.entries(m.groups).find(([, v]) => v !== undefined);
        if (!entry) continue;
        const [cls, value] = entry;

        out += escapeHtml(code.slice(last, m.index));
        out += `<span class="tok-${cls}">${escapeHtml(value)}</span>`;
        last = m.index + m[0].length;
    }

    return out + escapeHtml(code.slice(last));
}

function highlightAll(root = document) {
    if (!SHELL_RE) return;

    root.querySelectorAll('code[data-lang="shell"]').forEach(el => {
        // textContent, never innerHTML: whatever is in the DOM is the source of
        // truth, and re-running this is idempotent.
        el.innerHTML = highlightShell(el.textContent);
    });
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = { highlightShell };  // for the test harness
} else {
    highlightAll();
}
