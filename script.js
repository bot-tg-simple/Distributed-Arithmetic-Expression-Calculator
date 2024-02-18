const form = document.getElementById("expressionForm");
const input_express = document.getElementById("expressInput");
const results = document.querySelector(".results");

form.addEventListener("submit", function(e) {
    e.preventDefault(); // Предотвращаем стандартное действие браузера

    let expression = input_express.value;

    const express = input_express.value;

    if (expression === '') {
        return;
    }

    // Кодирование всей строки expression перед отправкой
    expression = encodeURIComponent(expression);
    
    fetch('http://localhost:8080/expression', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
            expression: expression,
        }),
    })
    .then((response) => response.json())
    .then((data) => {
        console.log('Expression ID:', data.id);
        results.innerHTML += `<div class="result success">
            <span class="status-icon"></span>
            <span class="expression">${express}  (ID: ${data.id})</span>
        </div>`;
    })
    .catch((error) => {
        alert(`Error: ${error}\n
        Скорее всего вы просто не запустили сервер`);
        console.error('Error:', error);
    });

    // Очистка поля ввода после отправки
    input_express.value = '';
});
