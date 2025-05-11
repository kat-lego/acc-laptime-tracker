import { formatMilliseconds, formatUnixToLocalDateTime } from '@/utilities/date'
import { Lap, Session } from '../models/session'
import { getCarName } from '@/utilities/kunosCarMap'

interface SessionInfoPaneParams {
  session: Session
}

export default function SessionInfoPane({ session }: SessionInfoPaneParams) {

  return (
    <>
      <SessionHero session={session} />
      <SessionLapTable session={session} />
    </>
  )
}


interface SessionHeroParams {
  session: Session
}

function SessionHero({ session }: SessionHeroParams) {

  return (
    <div className="hero">
      <div className="hero-content text-center">
        <div className="">
          <h2 className="text-4xl font-bold">{formatUnixToLocalDateTime(session.startTime)}</h2>
          <h1 className="text-3xl font-bold text-blue-500">{session.track}</h1>
          <h2 className="text-3xl font-bold text-pink-500">{getCarName(session.carModel)}</h2>
          <p className="py-6">
            Best Time: {formatMilliseconds(session.bestLapTime)} |
            Previous Time: {formatMilliseconds(session.previousLapTime)}
          </p>
          {session.isActive && <div className="badge badge-success">active</div>}
          {!session.isActive && <div className="badge badge-error">inactive</div>}
        </div>
      </div>
    </div>
  )
}


interface SessionLapTableParams {
  session: Session
}

function SessionLapTable({ session }: SessionLapTableParams) {

  return (
    <div className="overflow-x-auto">
      <table className="table">
        <thead>
          <tr>
            <th>Lap #</th>
            <th>Time</th>
            <th>Time Delta</th>
            <th>Validity</th>
            {Array.from({ length: session.numberOfSectors }, (_, i) => (
              <>
                <th>Sector {i + 1}</th >
              </>
            ))}
          </tr>
        </thead>
        <tbody>
          {session.laps.map((lap, _) => (
            <>
              <SessionLapTableRow lap={lap} />
            </>
          ))}
        </tbody>
      </table>
    </div >
  )
}


interface SessionLapTableRowParams {
  lap: Lap
}

function SessionLapTableRow({ lap }: SessionLapTableRowParams) {
  let selectionClass = lap.isActive ? "bg-base-200" : ""
  return (
    <tr className={selectionClass}>
      <th>{lap.lapNumber}</th>
      <td>{formatMilliseconds(lap.lapTime)}</td>
      <td>{formatMilliseconds(lap.lapTime)}</td>
      <td>{lap.isValid ? "valid" : "invalid"}</td>
      {lap.lapSectors.map((s, _) => (
        <>
          <td>{formatMilliseconds(s.sectorTime)}</td>
        </>
      ))}
    </tr>
  )
}
