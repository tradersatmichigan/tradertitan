import React, { useContext, useState } from "react";
import { StateContext } from "./Game";

const Make = () => {
  const state = useContext(StateContext);
  const [bid, setBid] = useState<number>();

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (bid === undefined) {
      return;
    }
    if (state.room.username !== "" && bid >= state.room.width) {
      alert("Bid must be narrower than current bid");
      return;
    }
    const formData = new FormData();
    formData.append("value", String(bid));
    fetch("/make", {
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
      <p>
        Narrowest bid:{" "}
        {state.room.username === ""
          ? "none"
          : `${state.room.width} by ${state.room.username}`}
      </p>
      <form onSubmit={handleSubmit}>
        <span>Enter bid width: </span>
        <input
          type="number"
          name="bid"
          placeholder="bid"
          value={bid === undefined ? "" : String(bid)}
          onChange={(e) => setBid(Number(e.target.value))}
          min="0"
          required
        />
        <button type="submit">Submit</button>
      </form>
    </>
  );
};

export default Make;
