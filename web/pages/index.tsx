import SessionInfoPane from "@/components/SessionInfoPane";
import { Session as AccSession } from "@/models/session";
import { formatUnixToLocalDateTime } from "@/utilities/date";
import { getCarName } from "@/utilities/kunosCarMap";
import { useEffect, useState } from "react";

interface AccSessionResponse {
  sessions: AccSession[]
  total: number
}

export default function Home() {

  const [sessions, setSessions] = useState<AccSession[]>([])
  const [selectedSession, setSelectedSession] = useState<number>(-1)

  useEffect(() => {
    const fetchSessions = async () => {
      const BASE_URL = 'http://localhost:8080'
      const response = await fetch(`${BASE_URL}/api/sessions`)
      const sessions = await response.json() as AccSessionResponse

      setSessions(sessions.sessions)
      if (sessions.sessions.length > 0) {
        setSelectedSession(0)
      }
    }

    fetchSessions()
  }, [])

  const handleSelectSession = (index: number) => {
    setSelectedSession(index)
  }

  return (
    <div
      className="justify-items-center min-h-screen gap-15
      font-[family-name:var(--font-geist-sans)] pb-16 overflow-x-hidden">

      <header className="sticky top-0 z-50 navbar bg-app-black h-fit">

        <div className="navbar-start">
          <div className="dropdown">
            <div tabIndex={0} role="button" className="btn btn-ghost btn-circle">
              <DropDownIcon />
            </div>
            <ul
              tabIndex={0}
              className="list dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 shadow">
              {sessions.map((s, i) => (
                <li className="list-row hover:bg-blue-600" key={s.id} onClick={() => handleSelectSession(i)}>
                  <div className="uppercase font-semibold opacity">
                    <div>{s.track}</div>
                    <div className="uppercase font-semibold opacity-60">
                      {getCarName(s.carModel)}
                    </div>
                    <div className="text-xs uppercase font-semibold opacity-80">
                      {formatUnixToLocalDateTime(s.startTime)}
                    </div>
                  </div>
                </li>
              ))}

              {/* <div className="join grid grid-cols-2"> */}
              {/*   <button className="join-item btn ">load prev</button> */}
              {/*   <button className="join-item btn">load next</button> */}
              {/* </div> */}
            </ul>
          </div>
        </div>

        <div className="navbar-end">
          <a className="btn btn-ghost text-xl">
            ACC LAPTIME TRACKER
          </a>
        </div>

      </header>

      <main className="relative flex flex-col gap-8
          max-w-screen-lg w-full px-6 xs:px-1">
        {sessions.length > 0 && selectedSession == -1
          && <span className="text-center">select a session</span>}
        {sessions.length > 0 && selectedSession >= 0
          && <SessionInfoPane session={sessions[selectedSession]} />}
      </main>
    </div>
  );
}


function DropDownIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className="h-5 w-5"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M4 6h16M4 12h16M4 18h7" />
    </svg>
  )
}
