document.addEventListener('DOMContentLoaded', () => {
    const mkdirButton = document.querySelector('.mkdir');
    if (mkdirButton) {
        mkdirButton.addEventListener('click', showMkdirModal);
    }
});

function showMkdirModal() {
    const modal = document.createElement('div');
    modal.className = 'action-modal';
    modal.innerHTML = `
        <div class="modal-content">
            <h3>Create New Folder</h3>
            <input type="text" id="new-folder-input" placeholder="Folder name" />
            <div class="button-group">
                <button class="action-confirm">Create</button>
                <button class="action-cancel">Cancel</button>
            </div>
        </div>
    `;
    document.body.appendChild(modal);

    const input = modal.querySelector('#new-folder-input');
    const confirmBtn = modal.querySelector('.action-confirm');
    const cancelBtn = modal.querySelector('.action-cancel');

    input.focus();

    function closeModal() {
        document.body.removeChild(modal);
    }
    cancelBtn.addEventListener('click', closeModal);
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });

    confirmBtn.addEventListener('click', () => {
        const folderName = input.value.trim();
        if (folderName !== '') {
            createFolder(folderName);
        }
        closeModal();
    });

    modal.style.display = 'block';
}


async function createFolder(folderName) {
    const currentPath = window.location.pathname.replace(/^\/drive\/?/, '');
    const folderPath = currentPath + folderName;

    try {
        const response = await fetch('/api', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'mkdir',
                path: folderPath
            })
        });
        const data = await response.json();
        if (data.success) {
            addFolderToGrid(folderName);
        } else {
            throw new Error(data.message);
        }
    } catch (error) {
        console.error('Error creating folder:', error);
        alert(`Error creating folder: ${error.message}`);
    }
}

function addFolderToGrid(folderName) {
    const grid = document.querySelector('.folders .grid');
    if (!grid) return;

    const folderItem = document.createElement('div');
    folderItem.className = 'item folder';

    const link = document.createElement('a');
    link.href = folderName + '/';

    const iconNameDiv = document.createElement('div');
    iconNameDiv.className = 'icon-name';

    const icon = document.createElement('i');
    icon.className = 'fas fa-folder';

    const span = document.createElement('span');
    span.title = folderName;
    span.textContent = folderName;

    iconNameDiv.appendChild(icon);
    iconNameDiv.appendChild(span);
    link.appendChild(iconNameDiv);
    folderItem.appendChild(link);

    grid.appendChild(folderItem);
}