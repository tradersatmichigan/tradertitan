import React, { useContext } from "react";
import { StateContext } from "./Game";

const Leaderboard = () => {
  const state = useContext(StateContext);

  if (state === undefined || state?.room.ranks === null) {
    return;
  }

  return (
    <div className="bg-white shadow-lg rounded-2xl p-6 max-w-sm w-full">
      <h3 className="text-xl font-semibold text-center text-gray-800">
        Room Leaderboard
      </h3>

      <ul className="mt-4 space-y-2 text-gray-700">
        {state.room.ranks.map(({ rank, username }) => (
          <li
            key={username}
            className="flex justify-between px-4 py-2 bg-gray-100 rounded-lg"
          >
            <span className="font-medium">{rank}.</span>
            <span>{username}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Leaderboard;
