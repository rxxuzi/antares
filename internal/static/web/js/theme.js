// Elements
const menuToggle = document.getElementById('menu-toggle');
const themeToggle = document.querySelector('.theme-toggle');
const sidebar = document.getElementById('sidebar-wrapper');
const content = document.getElementById('content-wrapper');

document.addEventListener('DOMContentLoaded', function () {
    let isMobile = window.innerWidth <= 768;

    function toggleSidebar() {
        if (isMobile) {
            sidebar.classList.toggle('mobile-active');
            content.classList.toggle('mobile-pushed');
        } else {
            document.documentElement.classList.toggle('sidebar-closed');
        }
        updateMenuToggleIcon();
        const isClosed = isMobile ? !sidebar.classList.contains('mobile-active') : document.documentElement.classList.contains('sidebar-closed');
        localStorage.setItem('sidebarClosed', isClosed ? 'true' : 'false');
    }

    function updateMenuToggleIcon() {
        const isClosed = isMobile
            ? !sidebar.classList.contains('mobile-active')
            : document.documentElement.classList.contains('sidebar-closed');
        menuToggle.innerHTML = isClosed ? '<i class="fas fa-bars"></i>' : '<i class="fas fa-times"></i>';
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

    function updateLayout() {
        isMobile = window.innerWidth <= 768;
        if (isMobile) {
            sidebar.classList.add('mobile');
            content.classList.add('mobile');
            document.documentElement.classList.remove('sidebar-closed');
            const sidebarClosed = localStorage.getItem('sidebarClosed') === 'true';
            sidebar.classList.toggle('mobile-active', !sidebarClosed);
            content.classList.toggle('mobile-pushed', !sidebarClosed);
        } else {
            sidebar.classList.remove('mobile', 'mobile-active');
            content.classList.remove('mobile', 'mobile-pushed');
            const sidebarClosed = localStorage.getItem('sidebarClosed') === 'true';
            document.documentElement.classList.toggle('sidebar-closed', sidebarClosed);
        }
        updateMenuToggleIcon();
    }

    function setDynamicSizes() {
        const vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty('--vh', `${vh}px`);

        const itemWidth = isMobile ? Math.floor((window.innerWidth - 30) / 2) : 200;
        document.documentElement.style.setProperty('--item-width', `${itemWidth}px`);
        document.documentElement.style.setProperty('--folder-height', `${Math.floor(itemWidth / 4)}px`);
        document.documentElement.style.setProperty('--file-height', `${itemWidth}px`);
    }

    function setVH() {
        let vh = window.innerHeight * 0.01;
        document.documentElement.style.setProperty('--vh', `${vh}px`);
    }

    // Initialize theme
    const savedTheme = localStorage.getItem('darkMode');
    const isDarkMode = savedTheme === 'dark';
    if (isDarkMode) document.documentElement.classList.add('dark-mode');
    updateThemeToggle(isDarkMode);

    const sidebarClosed = localStorage.getItem('sidebarClosed') === 'true';
    if (isMobile) {
        sidebar.classList.toggle('mobile-active', !sidebarClosed);
        content.classList.toggle('mobile-pushed', !sidebarClosed);
    } else {
        document.documentElement.classList.toggle('sidebar-closed', sidebarClosed);
    }
    updateMenuToggleIcon();

    // Event listeners
    menuToggle.addEventListener('click', toggleSidebar);
    themeToggle.addEventListener('click', toggleTheme);
    window.addEventListener('resize', () => {
        updateLayout();
        setDynamicSizes();
        setVH();
    });

    // Initial setup
    updateLayout();
    setDynamicSizes();
    setVH();
});