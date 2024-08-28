document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.querySelector('.search');
    const searchModal = document.getElementById('search-modal');
    const searchInput = document.getElementById('search-input');
    const searchSubmit = document.getElementById('search-submit');

    searchButton.addEventListener('click', (e) => {
        e.preventDefault();
        searchModal.style.display = 'flex';
        searchInput.focus();
    });

    searchSubmit.addEventListener('click', performSearch);
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            performSearch();
        }
    });

    function performSearch() {
        const query = searchInput.value.trim();
        if (query) {
            window.location.href = `/search?q=${encodeURIComponent(query)}`;
        }
    }

    searchModal.addEventListener('click', (e) => {
        if (e.target === searchModal) {
            searchModal.style.display = 'none';
        }
    });
});
