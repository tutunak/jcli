// Terminal typewriter animation
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
const cursor = document.querySelector('.cursor');
let currentSequence = 0;
let currentChar = 0;
let isTyping = false;

function typeWriter() {
    if (currentSequence >= sequences.length) {
        // Reset and loop
        setTimeout(() => {
            terminalOutput.innerHTML = '';
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
        terminalOutput.appendChild(span);
        currentChar++;

        const charDelay = seq.delay || 30;
        setTimeout(typeCharacter, charDelay);
    } else {
        currentSequence++;
        currentChar = 0;
        setTimeout(typeWriter, 100);
    }
}

// Copy to clipboard functionality
function copyInstall() {
    const cmd = document.getElementById('install-cmd').textContent;
    navigator.clipboard.writeText(cmd).then(() => {
        const btn = document.querySelector('.copy-btn');
        btn.classList.add('copied');
        setTimeout(() => btn.classList.remove('copied'), 2000);
    });
}

// Tab switching
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));

        tab.classList.add('active');
        document.getElementById(tab.dataset.tab).classList.add('active');
    });
});

// Start animation when page loads
document.addEventListener('DOMContentLoaded', () => {
    setTimeout(typeWriter, 500);
});
