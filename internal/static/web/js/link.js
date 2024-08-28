document.addEventListener("DOMContentLoaded", () => {
    handleParentDirLink();
})

function handleParentDirLink() {
    const parentDirLink = document.getElementById('parent-dir-link');
    const currentPath = window.location.pathname;
    if (currentPath === '/drive/' || currentPath === '/drive') {
        parentDirLink.style.display = 'none';
    } else {
        parentDirLink.style.display = 'block';
    }
}