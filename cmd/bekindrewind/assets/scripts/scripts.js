// Initialize the class change handler on DOMContentLoaded and htmx:afterSwap
document.addEventListener('DOMContentLoaded', handleVerifiedInputClassChange);
document.body.addEventListener('htmx:afterSwap', handleVerifiedInputClassChange);
document.body.addEventListener('htmx:afterSwap', handleFileImport);

// Invalid input and buttons
function handleVerifiedInputClassChange() {
    const verifiedInputs = document.querySelectorAll('.verified-input');
    const verifiedButtons = document.querySelectorAll('.verified-button');

    verifiedInputs.forEach(input => {
        if (input.classList.contains('invalid-input')) {
            verifiedButtons.forEach(button => {
                button.setAttribute('disabled', 'disabled');
            });
        }
        const observer = new MutationObserver(() => {
            let anyDanger = false;
            verifiedInputs.forEach(input => {
                if (input.classList.contains('invalid-input')) {
                    anyDanger = true;
                }
            });

            verifiedButtons.forEach(button => {
                if (anyDanger) {
                    button.setAttribute('disabled', 'disabled');
                } else {
                    button.removeAttribute('disabled');
                }
            });
        });

        observer.observe(input, { attributes: true, attributeFilter: ['class'] });
    });
}

// Toggle dark theme
const body = document.body;
const themeToggle = document.getElementById('switch-theme');

const savedTheme = localStorage.getItem('theme');
if (savedTheme === 'dark') {
    body.classList.add('dark');
}

themeToggle.addEventListener('click', () => {
    const isDark = body.classList.toggle('dark');
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
});

//File import-export
function handleFileImport() {
    const fileInput = document.getElementById('file-upload-input');
    const openButton = document.getElementById('open-file-dialog');
    const fileNameDisplay = document.getElementById('selected-file-name');

    openButton.addEventListener('click', () => {
        fileInput.click();
    });

    fileInput.addEventListener('change', () => {
        if (fileInput.files.length > 0) {
            fileNameDisplay.textContent = fileInput.files[0].name;
        } else {
            fileNameDisplay.textContent = 'No file selected';
        }
    });
}

// Handle file download for CSV files
document.addEventListener('htmx:afterRequest', function (evt) {
    const xhr = evt.detail.xhr;
    if (xhr.getResponseHeader('Content-Disposition') === 'attachment' && typeof xhr.getResponseHeader('HX-Download') === 'string') {
        const filename = xhr.getResponseHeader('HX-Download');
        console.log(xhr.response);
        const blob = new Blob([xhr.response], { type: 'text/csv' });
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
});