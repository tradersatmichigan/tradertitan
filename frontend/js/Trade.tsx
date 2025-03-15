import React, { useContext, useState } from "react";
import { StateContext } from "./Game";
import { Side } from "./types";

const Trade = () => {
  const state = useContext(StateContext);
  const [side, setSide] = useState<Side>(Side.Long);

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("value", String(side));
    fetch("/trade", {
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

  const buy_text = `BUY @ ${state.room.center + state.room.width / 2}`;
  const sell_text = `SELL @ ${state.room.center - state.room.width / 2}`;

  return (
    <div className="bg-white shadow-lg rounded-2xl p-8 w-96">
      <h2 className="text-2xl font-semibold text-center text-gray-800">
        Place Your Trade
      </h2>

      <div className="mt-4 text-gray-700">
        <p className="text-lg">
          Market: <span className="font-semibold">{state.market}</span>
        </p>
      </div>

      <form className="mt-6" onSubmit={handleSubmit}>
        <span className="block text-gray-600 font-medium mb-2">
          Enter trade:
        </span>
        <div className="flex gap-4 justify-center">
          <label className="flex items-center space-x-2">
            <input
              type="radio"
              value="long"
              name="side"
              checked={side === Side.Long}
              onChange={(e) =>
                setSide(e.target.value === "long" ? Side.Long : Side.Short)
              }
              className="w-5 h-5 text-blue-500"
            />
            <span className="text-gray-700">{buy_text}</span>
          </label>

          <label className="flex items-center space-x-2">
            <input
              type="radio"
              value="short"
              name="side"
              checked={side === Side.Short}
              onChange={(e) =>
                setSide(e.target.value === "short" ? Side.Short : Side.Long)
              }
              className="w-5 h-5 text-red-500"
            />
            <span className="text-gray-700">{sell_text}</span>
          </label>
        </div>

        <button
          type="submit"
          className="w-full mt-4 bg-blue-500 text-white font-semibold py-2 rounded-lg hover:bg-blue-600 transition"
        >
          Submit
        </button>
      </form>
    </div>
  );
};

export default Trade;
