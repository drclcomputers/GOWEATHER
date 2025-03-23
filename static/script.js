function checkUrlForWord(word) {
    const currentUrl = window.location.href;
    return currentUrl.includes(word);
  }
  
  window.onload = function () {
    document.querySelector('form').reset();
    const word = "query?";
    if (checkUrlForWord(word)) {
      document.getElementById('weather-info').style.display = 'block';
      document.getElementById('forecast').style.display = 'flex';
    }
  };