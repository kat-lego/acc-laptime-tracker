document.addEventListener('DOMContentLoaded', async () => {
  try {
    console.log(`Fetching sessions from ${trackerApiBaseUrl}/api/sessions`)
    const response = await fetch(`${trackerApiBaseUrl}/api/sessions`);
    if (!response.ok) throw new Error('Failed to fetch sessions');

    window.sessions = (await response.json()).sessions;

    document.dispatchEvent(new CustomEvent('sessions:ready', { detail: window.sessions }));

    const sidebarEl = document.getElementById('sessions-list');
    if (sidebarEl && Array.isArray(window.sessions)) {
      sidebarEl.innerHTML = '';
      window.sessions.forEach((session, _) => {
        const li = document.createElement('li');

        const a = document.createElement('a');
        a.href = "#";
        const car = getCarName(session.carModel);
        const track = session.track.toUpperCase();
        a.textContent = `${car} @ ${track}`;

        const p = document.createElement('p');
        p.textContent = formatUnixToLocalDateTime(session.startTime);

        a.addEventListener('click', (e) => {
          e.preventDefault();
          document.dispatchEvent(new CustomEvent('session:selected', {
            detail: session
          }));
        });

        li.appendChild(a);
        li.appendChild(p);
        sidebarEl.appendChild(li);
      });
    }

  } catch (err) {
    console.error('Session fetch error:', err);
  }
});

document.addEventListener('session:selected', (e) => {
  const session = e.detail;
  const laptimes = session.laps || [];

  const table = document.getElementById('laptimes-table');
  const thead = table.querySelector('thead');
  const tbody = table.querySelector('tbody');
  const noData = document.getElementById('no-data');

  const sessionTitle = document.getElementById('session-title');
  const car = getCarName(session.carModel);
  const track = session.track.toUpperCase();
  sessionTitle.textContent = `${car} @ ${track}`;

  const sessionDetails = document.getElementById('session-details');
  sessionDetails.textContent = `${session.sessionType} | ${formatUnixToLocalDateTime(session.startTime)}`;
  sessionDetails.style.display = '';

  // Clear previous content
  thead.innerHTML = '';
  tbody.innerHTML = '';

  if (laptimes.length > 0) {
    const numSectors = session.numberOfSectors;

    // Build table header
    const headerRow = document.createElement('tr');

    const lapHeader = document.createElement('th');
    lapHeader.textContent = 'Lap';
    headerRow.appendChild(lapHeader);

    for (let i = 0; i < numSectors; i++) {
      const sectorHeader = document.createElement('th');
      sectorHeader.textContent = `Sector ${i + 1}`;
      headerRow.appendChild(sectorHeader);
    }

    const lapTimeHeader = document.createElement('th');
    lapTimeHeader.textContent = 'Time Delta';
    headerRow.appendChild(lapTimeHeader);

    thead.appendChild(headerRow);

    // Build table rows
    laptimes.forEach((lap, index) => {
      const row = document.createElement('tr');

      // Add class for active lap
      if (lap.isActive) {
        row.classList.add('active-lap');
      }

      // Add class for invalid lap
      if (!lap.isValid) {
        row.classList.add('invalid-lap');
      }

      // Add class for best lap
      if (index === session.bestLap) {
        row.classList.add('best-lap');
      }

      // Lap number cell
      const lapNumberCell = document.createElement('td');
      lapNumberCell.textContent = lap.lapNumber;
      row.appendChild(lapNumberCell);

      // Sector times
      if (Array.isArray(lap.lapSectors)) {
        lap.lapSectors.forEach((sector) => {
          const sectorCell = document.createElement('td');
          sectorCell.textContent = formatMilliseconds(sector.sectorTime);
          row.appendChild(sectorCell);
        });

        // Pad missing sector cells
        const missingSectors = numSectors - lap.lapSectors.length;
        for (let i = 0; i < missingSectors; i++) {
          const emptyCell = document.createElement('td');
          emptyCell.textContent = '-';
          row.appendChild(emptyCell);
        }
      }

      // Lap time cell
      const lapTimeCell = document.createElement('td');
      lapTimeCell.textContent = formatMilliseconds(lap.lapDelta);
      row.appendChild(lapTimeCell);

      tbody.appendChild(row);
    });

    table.style.display = '';
    noData.style.display = 'none';
  } else {
    table.style.display = 'none';
    noData.style.display = '';
  }
});

function formatUnixToLocalDateTime(unixTimestamp) {
  const date = new Date(unixTimestamp * 1000);

  const day = date.getDate();
  const monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ];
  const month = monthNames[date.getMonth()];
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');

  return `${day} ${month} ${year} ${hours}:${minutes}`;
}

function formatMilliseconds(ms) {
  if (ms >= 3600000) {
    return '--:--:---';
  }

  const minutes = Math.floor(ms / (1000 * 60));
  const seconds = Math.floor((ms % (1000 * 60)) / 1000);
  const milliseconds = ms % 1000;

  const pad = (num, size) => num.toString().padStart(size, '0');

  return `${pad(minutes, 2)}:${pad(seconds, 2)}:${pad(milliseconds, 3)}`;
}


function getCarName(kunosId) {
  const kunosCarMap = {
    // GT3 - 2013 to 2018
    amr_v12_vantage_gt3: "Aston Martin Vantage V12 GT3 2013",
    audi_r8_lms: "Audi R8 LMS 2015",
    bentley_continental_gt3_2016: "Bentley Continental GT3 2015",
    bentley_continental_gt3_2018: "Bentley Continental GT3 2018",
    bmw_m6_gt3: "BMW M6 GT3 2017",
    jaguar_g3: "Emil Frey Jaguar G3 2012",
    ferrari_488_gt3: "Ferrari 488 GT3 2018",
    honda_nsx_gt3: "Honda NSX GT3 2017",
    lamborghini_gallardo_rex: "Lamborghini Gallardo G3 Reiter 2017",
    lamborghini_huracan_gt3: "Lamborghini Huracan GT3 2015",
    lamborghini_huracan_st: "Lamborghini Huracan ST 2015",
    lexus_rc_f_gt3: "Lexus RCF GT3 2016",
    mclaren_650s_gt3: "McLaren 650S GT3 2015",
    mercedes_amg_gt3: "Mercedes AMG GT3 2015",
    nissan_gt_r_gt3_2017: "Nissan GTR Nismo GT3 2015",
    nissan_gt_r_gt3_2018: "Nissan GTR Nismo GT3 2018",
    porsche_991_gt3_r: "Porsche 991 GT3 R 2018",
    porsche_991ii_gt3_cup: "Porsche9 91 II GT3 Cup 2017",

    // GT3 - 2019
    amr_v8_vantage_gt3: "Aston Martin V8 Vantage GT3 2019",
    audi_r8_lms_evo: "Audi R8 LMS Evo 2019",
    honda_nsx_gt3_evo: "Honda NSX GT3 Evo 2019",
    lamborghini_huracan_gt3_evo: "Lamborghini Huracan GT3 EVO 2019",
    mclaren_720s_gt3: "McLaren 720S GT3 2019",
    porsche_991ii_gt3_r: "Porsche 911 II GT3 R 2019",

    // GT4
    alpine_a110_gt4: "Alpine A110 GT4 2018",
    amr_v8_vantage_gt4: "Aston Martin Vantage AMR GT4 2018",
    audi_r8_gt4: "Audi R8 LMS GT4 2016",
    bmw_m4_gt4: "BMW M4 GT4 2018",
    chevrolet_camaro_gt4r: "Chevrolet Camaro GT4 R 2017",
    ginetta_g55_gt4: "Ginetta G55 GT4 2012",
    ktm_xbow_gt4: "Ktm Xbow GT4 2016",
    maserati_mc_gt4: "Maserati Gran Turismo MC GT4 2016",
    mclaren_570s_gt4: "McLaren 570s GT4 2016",
    mercedes_amg_gt4: "Mercedes AMG GT4 2016",
    porsche_718_cayman_gt4_mr: "Porsche 718 Cayman GT4 MR 2019",

    // GT3 – 2020
    ferrari_488_gt3_evo: "Ferrari 488 GT3 Evo 2020",
    mercedes_amg_gt3_evo: "Mercedes AMG GT3 Evo 2020",

    // GT3 – 2021
    bmw_m4_gt3: "BMW M4 vGT3 2021",

    // Challengers Pack – 2022
    audi_r8_lms_evo_ii: "Audi R8 LMS Evo II 2022",
    bmw_m2_cs_racing: "BMW M2 Cup 2020",
    ferrari_488_challenge_evo: "Ferrari 488 Challenge Evo 2020",
    lamborghini_huracan_st_evo2: "Lamborghini Huracan ST Evo2 2021",
    porsche_992_gt3_cup: "Porsche 992 GT3 Cup 2021"
  };

  return kunosCarMap[kunosId] ?? kunosId;
}
