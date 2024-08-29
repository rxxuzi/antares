// Elements
const uploadListItem = document.querySelector('.file-upload');
const fileInput = document.createElement('input');

document.addEventListener('DOMContentLoaded', function() {
    fileInput.type = 'file';
    fileInput.style.display = 'none';
    document.body.appendChild(fileInput);

    uploadListItem.addEventListener('click', function(e) {
        e.preventDefault();
        fileInput.click();
    });

    fileInput.addEventListener('change', function(e) {
        if (this.files && this.files[0]) {
            const selectedFile = this.files[0];
            uploadFile(selectedFile);
        }
    });

    function uploadFile(file) {
        const formData = new FormData();
        formData.append('file', file);

        const xhr = new XMLHttpRequest();
        xhr.open('POST', window.location.pathname, true);

        xhr.onload = function() {
            if (xhr.status === 200) {
                window.location.reload();
            } else {
                alert('Upload failed: ' + xhr.statusText);
            }
        };

        xhr.send(formData);
    }
});