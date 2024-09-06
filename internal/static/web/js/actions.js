// Elements
const actionModal = document.getElementById('action-modal');
const folderActionModal = document.getElementById('folder-action-modal');

let currentFileInfo = null;
let currentOpenModal = null;

document.addEventListener('DOMContentLoaded', () => {
    setupItemListeners();
    setupActionListeners();
    setupGlobalListeners();
});


// Functions
function setupItemListeners() {
    document.querySelectorAll('.item').forEach(item => {
        item.addEventListener('contextmenu', handleRightClick);
    });
}

function setupActionListeners() {
    // -> file
    document.getElementById('action-download').addEventListener('click', handleDownload);
    document.getElementById('action-rename').addEventListener('click', handleRename);
    document.getElementById('action-copy').addEventListener('click', handleCopy);
    document.getElementById('action-move').addEventListener('click', handleMove);
    document.getElementById('action-info').addEventListener('click', handleInfo);
    document.getElementById('action-delete').addEventListener('click', handleDelete);

    // -> folder
    document.getElementById('folder-action-rename').addEventListener('click', handleRenameFolder);
    document.getElementById('folder-action-move').addEventListener('click', handleMoveFolder);
    document.getElementById('folder-action-delete').addEventListener('click', handleDeleteFolder);
}

function setupGlobalListeners() {
    document.addEventListener('click', (e) => {
        if (!e.target.closest('#action-modal') && !e.target.closest('#folder-action-modal')) {
            closeAllModals();
        }
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeAllModals();
        }
    });
}

function handleRightClick(e) {
    e.preventDefault();
    const item = e.currentTarget;
    const link = item.querySelector('a');

    currentFileInfo = {
        name: link.getAttribute('filename') || link.textContent.trim(),
        path: link.getAttribute('href'),
        type: item.classList.contains('folder') ? 'folder' : 'file',
        size: link.getAttribute('data-size'),
        modified: link.getAttribute('data-modified')
    };

    closeAllModals();

    if (currentFileInfo.type === 'folder') {
        showModal(folderActionModal, e.clientX, e.clientY);
    } else {
        showModal(actionModal, e.clientX, e.clientY);
    }
}

function showModal(modal, x, y) {
    modal.style.display = 'block';
    modal.style.left = `${x}px`;
    modal.style.top = `${y}px`;
    currentOpenModal = modal;
}

function closeAllModals() {
    actionModal.style.display = 'none';
    folderActionModal.style.display = 'none';
    currentOpenModal = null;
}

// Action Handlers
function handleDownload() {
    console.log('Downloading:', currentFileInfo.name);
    // TODO: 実際のダウンロード処理を実装
    closeAllModals();
}


function handleRename() {
    console.log('Renaming:', currentFileInfo.name);
    showRenameModal(currentFileInfo);
}

function handleRenameFolder() {
    console.log('Renaming folder:', currentFileInfo.name);
    showRenameModal(currentFileInfo);
}

function showRenameModal(fileInfo) {
    const renameModal = document.createElement('div');
    renameModal.className = 'action-modal';
    renameModal.innerHTML = `
        <div class="modal-content">
            <h3>Rename ${fileInfo.type === 'folder' ? 'Folder' : 'File'}</h3>
            <input type="text" id="new-name-input" value="${fileInfo.name}" />
            <div class="button-group">
                <button class="action-confirm">Rename</button>
                <button class="action-cancel">Cancel</button>
            </div>
        </div>
    `;
    document.body.appendChild(renameModal);

    const newNameInput = renameModal.querySelector('#new-name-input');
    const confirmButton = renameModal.querySelector('.action-confirm');
    const cancelButton = renameModal.querySelector('.action-cancel');

    newNameInput.focus();
    newNameInput.setSelectionRange(0, newNameInput.value.lastIndexOf('.'));

    function closeModal() {
        document.body.removeChild(renameModal);
    }

    confirmButton.addEventListener('click', () => {
        const newName = newNameInput.value;
        if (newName && newName !== fileInfo.name) {
            renameFile(fileInfo, newName);
        }
        closeModal();
    });

    cancelButton.addEventListener('click', closeModal);
    renameModal.addEventListener('click', (event) => {
        if (event.target === renameModal) {
            closeModal();
        }
    });

    renameModal.style.display = 'block';
}

async function renameFile(fileInfo, newName) {
    const currentPath = window.location.pathname.replace(PREFIX_DRIVE, '');
    const oldPath = currentPath + fileInfo.name;
    const newPath = currentPath + newName;

    console.log('Renaming:', oldPath, 'to', newPath);

    try {
        const response = await fetch('/api', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                file: fileInfo.type !== 'folder',
                type: 'rename',
                path: oldPath,
                dst: newPath
            })
        });

        const data = await response.json();
        console.log('Server response:', data);

        if (data.success) {
            alert(data.message);
            location.reload();
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('Error:', error);
        alert(`An error occurred while renaming: ${error.message}`);
    }
}

function handleCopy() {
    console.log('Copying:', currentFileInfo.name);
    // TODO: 実際のコピー処理を実装
    closeAllModals();
}

function handleMove() {
    console.log('Moving:', currentFileInfo.name);
    // TODO: 実際の移動処理を実装
    closeAllModals();
}

function handleInfo() {
    console.log('File Info:', currentFileInfo);
    showFileInfo(currentFileInfo);
    closeAllModals();
}

function handleDelete() {
    console.log('Deleting:', currentFileInfo.name);
    if (confirm(`Are you sure you want to delete "${currentFileInfo.name}"?`)) {
        // TODO: 実際の削除処理を実装
        console.log('File deleted:', currentFileInfo.name);
    }
    closeAllModals();
}

// Folder Action Handlers
function handleMoveFolder() {
    console.log('Moving folder:', currentFileInfo.name);
    // TODO: フォルダーの移動処理を実装
    closeAllModals();
}

function handleDeleteFolder() {
    console.log('Deleting folder:', currentFileInfo.name);
    if (confirm(`Are you sure you want to delete the folder "${currentFileInfo.name}" and all its contents?`)) {
        // TODO: フォルダーの削除処理を実装
        console.log('Folder deleted:', currentFileInfo.name);
    }
    closeAllModals();
}

function showFileInfo(fileInfo) {
    const infoModal = document.createElement('div');
    infoModal.className = 'modal';
    infoModal.innerHTML = `
        <div class="modal-content modal">
            <h3>File Information</h3>
            <p><strong>Name:</strong> ${fileInfo.name}</p>
            <p><strong>Type:</strong> ${fileInfo.type}</p>
            <p><strong>Size:</strong> ${fileInfo.size}</p>
            <p><strong>Modified:</strong> ${fileInfo.modified}</p>
            <p><strong>Path:</strong> ${fileInfo.path}</p>
            <button id="close-info">Close</button>
        </div>
    `;
    document.body.appendChild(infoModal);

    const closeButton = infoModal.querySelector('#close-info');
    closeButton.addEventListener('click', () => {
        document.body.removeChild(infoModal);
    });

    infoModal.style.display = 'block';
}