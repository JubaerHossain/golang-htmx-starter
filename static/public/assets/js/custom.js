document.addEventListener('DOMContentLoaded', function () {
    const toggleModeButton = document.getElementById('toggleModeButton');
    const body = document.body;

    // Check if a mode preference is stored
    const modePreference = localStorage.getItem('modePreference');
    if (modePreference) {
        body.classList.add(modePreference);
    }

    toggleModeButton.addEventListener('click', function () {
        // Toggle between dark and white modes by toggling the body class
        body.classList.toggle('dark-mode');
        body.classList.toggle('white-mode');

        // Store the mode preference in localStorage
        const currentMode = body.classList.contains('dark-mode') ? 'dark-mode' : 'white-mode';
        localStorage.setItem('modePreference', currentMode);
    });
});