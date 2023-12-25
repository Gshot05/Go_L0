document.getElementById("idForm").addEventListener("submit", function(event) {
    event.preventDefault(); 

    var idValue = document.getElementById("idInput").value;

    fetch("/get-info/" + idValue)
    .then(response => response.json())
    .then(data => {
        console.log("Ответ от сервера:", data);

        document.getElementById("result").innerHTML = `
            <p>Информация по ID ${idValue}:</p>
            <pre>${JSON.stringify(data, null, 2)}</pre>
        `;
    })
    .catch(error => {
        console.error("Ошибка:", error);
    });
});
