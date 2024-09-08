package global

const Page404 = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Antares - 404</title>
    <script>
        (function() {
            const darkMode = localStorage.getItem('darkMode');
            if (darkMode === 'dark') {
                document.documentElement.classList.add('dark-mode');
            }
        })();
    </script>
    <link rel="stylesheet" href="/web/css/styles.css">
    <link rel="stylesheet" href="/web/css/404.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <link rel="icon" type="image" href="/web/favicon.ico">
</head>
<body>
<div class="error-container">
    <h2 class="error-code">404</h2>
    <p class="error-message">Oops! Page Not Found</p>
    <a href="/" class="home-button">Return to Home</a>
</div>
</body>
</html>
`

const PageApi = `
<html>
<head>
	<title>Antares - API</title>
	<link rel="stylesheet" href="/web/css/1.css">
</head>
<body>
	<main>
		<h1>API Endpoint</h1>
		<p>This page is displayed when accessed from a browser.</p>
		<p>For detailed usage of this API, please refer to the <a href="https://github.com/rxxuzi/antares/blob/main/doc/api.md">api.md</a>.</p>
	</main>
</body>
</html>
`
