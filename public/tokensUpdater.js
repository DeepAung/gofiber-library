const minute = 60 * 1000;

updateTokens();
setInterval(async () => {
  await updateTokens();
}, 14 * minute);

async function updateTokens() {
  console.log("updateTokens");
  try {
    await fetch("/api/refresh", { method: "POST" });
  } catch (err) {
    console.log(err);
  }
}
