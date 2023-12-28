const minute = 60 * 1000;

setInterval(() => {
  fetch("/api/refresh", { method: "POST" })
    .catch((err) => {
      console.log(err);
    });
}, 14 * minute);
