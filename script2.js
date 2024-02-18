const form = document.getElementById("expressionForm");
const results = document.getElementById("computing-resources");

form.addEventListener("submit", function(e) {
    e.preventDefault(); // Предотвращаем стандартное поведение браузера
    
    fetch('http://localhost:8080/ping', {
        method: 'GET'
    })
    .then((response) => {
        if (response.status === 404) {
            results.innerHTML += `<article class="result success">
            <span class="status-icon"></span>
            <span class="expression">
            localhost
            <br>
            last ping -
            <br>
            the connection to the server is stable
            </span>
            </article>`;
        } else {
            return response.json();
        }
    })
    .then((data) => {
        results.innerHTML += `<article class="result error">
            <span class="status-icon"></span>
            <span class="expression">
            localhost
            <br>
            last ping -
            <br>
            ${data.error}
            </span>
            </article>`;
            alert('Сервер остановлен!');
    })
    .catch((error) => {
        if (error.response && error.response.status !== 404) {
        alert(`Error: ${error}\n
        Вероятно, вы просто не запустили сервер`);
        console.error('Error:', error);
        }
    });
    
    return false;
});
