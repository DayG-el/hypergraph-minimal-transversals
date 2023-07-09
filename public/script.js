window.onload = function () {
    var form = document.getElementById('inputForm');
    var inputData = document.getElementById('inputData');
    var resultContainer = document.getElementById('resultContainer');
    var resultOutput = document.getElementById('resultOutput');
  
    form.addEventListener('submit', function (event) {
      event.preventDefault();
  
      var input = inputData.value;
  
      var requestData = {
        data: input.split('\n').filter(function (line) {
          return line.trim() !== '';
        })
      };
  
      fetch('/process', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
      })
        .then(function (response) {
          if (!response.ok) {
            throw new Error('Request failed: ' + response.status);
          }
          return response.json();
        })
        .then(function (data) {
          var results = data.results.join('\n');
          resultOutput.textContent = results;
          resultContainer.style.display = 'block';
        })
        .catch(function (error) {
          console.error('Error:', error);
        });
    });
  };
  