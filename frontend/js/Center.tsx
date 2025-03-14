import React, { useContext, useState } from "react";
import { StateContext } from "./Game";

const Center = () => {
  const state = useContext(StateContext);
  const [center, setCenter] = useState<number>();

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (center === undefined) {
      return;
    }
    const formData = new FormData();
    formData.append("value", String(center));
    fetch("/center", {
      method: "POST",
      body: formData,
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(await response.text());
        }
        return response.text();
      })
      .catch((error) => console.error("Fetch error:", error));
  };

  return (
    <>
      <p>Market: {state.market}</p>
      <p>Your bid: {state.room.width}</p>
      <p>Current center: {state.room.center}</p>
      <form onSubmit={handleSubmit}>
        <span>Center market: </span>
        <input
          type="number"
          name="bid"
          placeholder="center"
          value={center === undefined ? "" : String(center)}
          onChange={(e) => setCenter(Number(e.target.value))}
          min="0"
          required
        />
        <button type="submit">Submit</button>
      </form>
    </>
  );
};

export default Center;
