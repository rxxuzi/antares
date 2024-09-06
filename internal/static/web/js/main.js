// Elements
const PREFIX_DRIVE = '/drive/';

const fileLinks = document.querySelectorAll('.file-link');
const preview = document.getElementById('preview');
const viewerContent = document.getElementById('viewer-content');
const closeBtn = document.querySelector('.close-button');
const downloadBtn = document.querySelector('.download-button');
const infoBtn = document.querySelector('.info-button');
const filename = document.querySelector('.filename');
const fileInfo = document.getElementById('file-info');
const filePath = document.getElementById('file-path');
const fileModified = document.getElementById('file-modified');
const fileSize = document.getElementById('file-size');
const fileType = document.getElementById('file-type');
const logo = document.getElementById('logo');
const deleteBtn = document.querySelector('.delete-button');

document.addEventListener('DOMContentLoaded', () => {
    let currentFile = null;
    let currentFileIndex = -1;
    // main.js の該当部分を以下のように修正します

    let touchStartX = 0;
    let touchStartY = 0;
    let touchEndX = 0;
    let touchEndY = 0;


    // Functions
    function openPreview(link) {
        const file = link.getAttribute('data-filename');
        if (!file) {
            console.error('Filename is missing');
            return;
        }
        currentFile = {
            name: file,
            path: link.href,
            type: link.closest('.item').classList[2],
            size: link.getAttribute('data-size'),
            modified: link.getAttribute('data-modified')
        };

        const fileExtension = file.split('.').pop().toLowerCase();
        const fileType = link.closest('.item').classList[2];

        filename.textContent = file;

        updatePreviewContent(fileType, fileExtension, link.href);
        preview.style.display = 'block';
        document.body.classList.add('preview-active');

        document.querySelector('.preview-header').style.display = 'flex';
        fileInfo.classList.add('hidden');
    }

    function updatePreviewContent(fileType, fileExtension, filePath) {
        switch (fileType) {
            case 'image':
                viewerContent.innerHTML = `<img src="${filePath}" alt="${currentFile.name}">`;
                break;
            case 'audio':
                viewerContent.innerHTML = `<audio controls src="${filePath}">Your browser does not support the audio element.</audio>`;
                break;
            case 'video':
                viewerContent.innerHTML = `<video controls><source src="${filePath}" type="video/${fileExtension}">Your browser does not support the video tag.</video>`;
                break;
            case 'pdf':
                viewerContent.innerHTML = `<iframe src="${filePath}" style="width: 100%; height: 100%;"></iframe>`;
                break;
            case 'text':
            case 'code':
                fetch(filePath)
                    .then(response => response.text())
                    .then(text => {
                        viewerContent.innerHTML = `<pre><code class="language-${fileExtension}">${escapeHtml(text)}</code></pre>`;
                    });
                break;
            default:
                viewerContent.innerHTML = `
                    <div class="binary-preview">
                        <img src="/web/assets/pv-binary.svg" alt="File icon" id="binary-img">
                        <p>Preview not available for this file type.</p>
                    </div>
                `;
        }
    }

    function closePreview() {
        preview.style.display = 'none';
        document.body.classList.remove('preview-active');
        viewerContent.innerHTML = '';

        document.querySelector('.preview-header').style.display = 'none';
        fileInfo.classList.add('hidden');
        currentFile = null;
        currentFileIndex = -1;
    }


    function escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    function showNextFile() {
        if (currentFileIndex < fileLinks.length - 1) {
            currentFileIndex++;
            openPreview(fileLinks[currentFileIndex]);
        }
    }

    function showPreviousFile() {
        if (currentFileIndex > 0) {
            currentFileIndex--;
            openPreview(fileLinks[currentFileIndex]);
        }
    }

    function handleSwipe() {
        const deltaX = touchEndX - touchStartX;
        const deltaY = touchEndY - touchStartY;

        // 縦方向のスワイプの方が大きい場合は、ファイル切り替えを行わない
        if (Math.abs(deltaY) > Math.abs(deltaX)) {
            return;
        }

        if (deltaX < -50) showNextFile();
        if (deltaX > 50) showPreviousFile();
    }

    function deleteFile() {
        if (currentFile) {
            const confirmDelete = confirm(`Are you sure you want to delete "${currentFile.name}"?`);
            if (confirmDelete) {
                const currentPath = window.location.pathname.replace(PREFIX_DRIVE, '');
                const filePath = currentPath + currentFile.name;

                console.log('Deleting file:', filePath);

                fetch('/api', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        type: 'rm',
                        path: filePath
                    })
                })
                    .then(response => response.json())
                    .then(data => {
                        console.log('Server response:', data);
                        if (data.success) {
                            alert(data.message);
                            closePreview();
                            location.reload();
                        } else {
                            throw new Error(data.message);
                        }
                    })
                    .catch(error => {
                        console.error('Error:', error);
                        alert(`An error occurred while deleting the file: ${error.message}`);
                    });
            }
        } else {
            console.error('No file selected for deletion');
        }
    }

    // Event Listeners
    viewerContent.addEventListener('touchstart', (e) => {
        touchStartX = e.changedTouches[0].screenX;
        touchStartY = e.changedTouches[0].screenY;
    });

    viewerContent.addEventListener('touchend', (e) => {
        touchEndX = e.changedTouches[0].screenX;
        touchEndY = e.changedTouches[0].screenY;
        handleSwipe();
    });

    fileLinks.forEach((link, index) => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            currentFileIndex = index;
            openPreview(link);
        });
    });
    closeBtn.addEventListener('click', closePreview);
    downloadBtn.addEventListener('click', () => {
        if (currentFile) {
            const link = document.createElement('a');
            link.href = currentFile.path;
            link.download = currentFile.name;
            link.click();
        }
    });
    infoBtn.addEventListener('click', () => {
        if (currentFile) {
            filePath.textContent = currentFile.path;
            fileModified.textContent = currentFile.modified;
            fileSize.textContent = currentFile.size;
            fileType.textContent = currentFile.type;
            fileInfo.classList.toggle('hidden');
        }
    });
    viewerContent.addEventListener('click', (e) => {
        if (e.target === viewerContent) closePreview();
    });
    logo.addEventListener('click', () => location.href = '/');
    document.addEventListener('keydown', (e) => {
        if (preview.style.display === 'block') {
            if (e.key === 'ArrowRight') {
                e.preventDefault();
                showNextFile();
            } else if (e.key === 'ArrowLeft') {
                e.preventDefault();
                showPreviousFile();
            }
        }
    });

    // イベントリスナーの追加
    deleteBtn.addEventListener('click', deleteFile);
    // Initialization
    const savedTheme = localStorage.getItem('darkMode');
    const isDarkMode = savedTheme === 'dark';
    if (isDarkMode) document.documentElement.classList.add('dark-mode');


});