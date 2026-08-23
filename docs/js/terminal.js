// Terminal typewriter animation.
//
// The transcript is already present in index.html as real text so that crawlers,
// language models and no-JS visitors can read it. This script only takes over
// when it can actually animate: it snapshots nothing, wipes the <pre>, and
// retypes the same content. Visitors who prefer reduced motion keep the static
// transcript untouched.
const sequences = [
    { text: '$ ', class: 'prompt', delay: 0 },
    { text: 'jcli issue select', class: 'command', delay: 50 },
    { text: '\n', delay: 400 },
    { text: '? Select an issue\n', class: 'output', delay: 30 },
    { text: '> ', class: 'selected', delay: 0 },
    { text: 'PROJ-123: Implement user authentication\n', class: 'selected', delay: 0 },
    { text: '  PROJ-124: Fix login redirect bug\n', class: 'output', delay: 0 },
    { text: '  PROJ-125: Add password reset flow\n', class: 'output', delay: 0 },
    { text: '\n', delay: 1500 },
    { text: '✓ Selected PROJ-123\n\n', class: 'selected', delay: 0 },
    { text: '$ ', class: 'prompt', delay: 500 },
    { text: 'git checkout -b $(jcli issue branch)', class: 'command', delay: 40 },
    { text: '\n', delay: 400 },
    { text: "Switched to a new branch '", class: 'output', delay: 0 },
    { text: 'PROJ-123-implement-user-authentication-847291', class: 'highlight', delay: 0 },
    { text: "'\n\n", class: 'output', delay: 0 },
    { text: '$ ', class: 'prompt', delay: 800 },
    { text: 'jcli issue current', class: 'command', delay: 50 },
    { text: '\n', delay: 400 },
    { text: 'PROJ-123', class: 'highlight', delay: 0 },
    { text: ': Implement user authentication\n', class: 'output', delay: 0 },
];

const terminalOutput = document.getElementById('terminal-output');
// The cursor lives inside the <pre> so it trails the typed text instead of
// dropping to its own line. Characters are always inserted before it.
const cursor = terminalOutput && terminalOutput.querySelector('.cursor');
let currentSequence = 0;
let currentChar = 0;

function typeWriter() {
    if (currentSequence >= sequences.length) {
        // Reset and loop
        setTimeout(() => {
            terminalOutput.replaceChildren(cursor);
            currentSequence = 0;
            currentChar = 0;
            typeWriter();
        }, 3000);
        return;
    }

    const seq = sequences[currentSequence];

    if (currentChar === 0 && seq.delay > 0) {
        setTimeout(() => {
            currentChar = 0;
            typeCharacter();
        }, seq.delay);
    } else {
        typeCharacter();
    }
}

function typeCharacter() {
    const seq = sequences[currentSequence];

    if (currentChar < seq.text.length) {
        const span = document.createElement('span');
        span.className = seq.class || '';
        span.textContent = seq.text[currentChar];
        terminalOutput.insertBefore(span, cursor);
        currentChar++;

        const charDelay = seq.delay || 30;
        setTimeout(typeCharacter, charDelay);
    } else {
        currentSequence++;
        currentChar = 0;
        setTimeout(typeWriter, 100);
    }
}

function startTerminal() {
    if (!terminalOutput || !cursor) return;
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) return;

    terminalOutput.replaceChildren(cursor);
    setTimeout(typeWriter, 500);
}

// Copy to clipboard functionality
function initCopyButton() {
    const btn = document.querySelector('.copy-btn');
    const cmd = document.getElementById('install-cmd');
    if (!btn || !cmd) return;

    btn.addEventListener('click', async () => {
        try {
            await navigator.clipboard.writeText(cmd.textContent.trim());
        } catch {
            return;
        }
        btn.classList.add('copied');
        setTimeout(() => btn.classList.remove('copied'), 2000);
    });
}

// Tab switching, wired to the ARIA tablist in index.html
function initTabs() {
    const tabs = Array.from(document.querySelectorAll('.tab'));
    if (tabs.length === 0) return;

    function select(tab, focus) {
        tabs.forEach(t => {
            const panel = document.getElementById(t.dataset.tab);
            const isActive = t === tab;

            t.classList.toggle('active', isActive);
            t.setAttribute('aria-selected', String(isActive));
            t.tabIndex = isActive ? 0 : -1;

            if (panel) {
                panel.classList.toggle('active', isActive);
                panel.hidden = !isActive;
            }
        });
        if (focus) tab.focus();
    }

    tabs.forEach((tab, i) => {
        tab.addEventListener('click', () => select(tab, false));
        tab.addEventListener('keydown', e => {
            const step = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : 0;
            if (step === 0) return;
            e.preventDefault();
            select(tabs[(i + step + tabs.length) % tabs.length], true);
        });
    });
}

initCopyButton();
initTabs();
startTerminal();
