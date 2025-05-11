const kunosCarMap: Record<string, string> = {
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

export function getCarName(kunosId: string): string {
  return kunosCarMap[kunosId] ?? kunosId;
}

