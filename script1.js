const input_express = document.getElementById("expressInput");
const info = document.querySelector(".expression");

input_express.addEventListener("keydown", function(e) {
    if (e.key === "Enter") {
        const expression = input_express.value;
        if (expression === '') {
            return;
        }
        
        fetch(`http://localhost:8080/expression/${expression}`, {
            method: 'GET'
        })
        .then((response) => response.json())
        .then((data) => {

            const expressionText = data.Expression.replace(/%2B/g, '+');

            info.innerHTML = ''; // Очищаем содержимое элемента перед добавлением новой информации
            info.innerHTML = `ID: ${data.ID}<br>
            Expression: ${expressionText}<br>
            Status: ${data.Status}<br>
            CreatedAt: ${data.CreatedAt}<br>
            UpdatedAt: ${data.UpdatedAt}<br>
            Result: ${data.Result}`;
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }
});
