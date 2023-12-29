const minute = 60 * 1000;

updateRefreshToken();
setInterval(() => {
  updateRefreshToken();
}, 14 * minute);

function updateRefreshToken() {
  console.log("updateRefreshToken");
  fetch("localhost:8080/api/refresh", { method: "POST" })
    .catch((err) => {
      console.log(err);
    });
}
