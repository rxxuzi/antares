document.addEventListener('DOMContentLoaded', () => {
    const searchButton = document.querySelector('.search');
    const searchModal = document.getElementById('search-modal');
    const searchInput = document.getElementById('search-input');
    const searchSubmit = document.getElementById('search-submit');
    const caseSensitiveToggle = document.getElementById('case-sensitive-toggle');
    const useRegexToggle = document.getElementById('use-regex-toggle');

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
            let searchUrl = '/search?';
            if (useRegexToggle.checked) searchUrl += 'r&';
            if (caseSensitiveToggle.checked) searchUrl += 'c&';
            searchUrl += `q=${encodeURIComponent(query)}`;
            window.location.href = searchUrl;
        }
    }

    searchModal.addEventListener('click', (e) => {
        if (e.target === searchModal) {
            searchModal.style.display = 'none';
        }
    });
});