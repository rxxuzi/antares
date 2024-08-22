document.addEventListener('DOMContentLoaded', function () {
    const menuToggle = document.getElementById('menu-toggle');
    const themeToggle = document.querySelector('.theme-toggle');

    function toggleSidebar() {
        const isClosed = document.documentElement.classList.toggle('sidebar-closed');
        localStorage.setItem('sidebarClosed', isClosed ? 'true' : 'false');
        updateMenuToggleIcon(isClosed);
    }

    function updateMenuToggleIcon(isClosed) {
        menuToggle.innerHTML = isClosed
            ? '<i class="fas fa-bars"></i>'
            : '<i class="fas fa-times"></i>';
    }

    function toggleTheme() {
        const isDarkMode = document.documentElement.classList.toggle('dark-mode');
        localStorage.setItem('darkMode', isDarkMode ? 'dark' : 'light');
        updateThemeToggle(isDarkMode);
    }

    function updateThemeToggle(isDarkMode) {
        const icon = themeToggle.querySelector('i');
        const text = themeToggle.querySelector('span');
        if (isDarkMode) {
            icon.classList.replace('fa-moon', 'fa-sun');
            text.textContent = 'Light';
        } else {
            icon.classList.replace('fa-sun', 'fa-moon');
            text.textContent = 'Dark';
        }
    }

    const savedTheme = localStorage.getItem('darkMode');
    const isDarkMode = savedTheme === 'dark';
    updateThemeToggle(isDarkMode);

    const sidebarClosed = localStorage.getItem('sidebarClosed') === 'true';
    updateMenuToggleIcon(sidebarClosed);

    menuToggle.addEventListener('click', toggleSidebar);
    themeToggle.addEventListener('click', toggleTheme);
});