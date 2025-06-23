package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("--Microservice Migrate Start--")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dbURL := os.Getenv("POSTGRES_URI")
	m, err := migrate.New(
		"file://files",
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed create migrate instance : %v", err)
	}

	// sdadsa
	err = m.Up()
	// err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully")

}


  procedure x_submit_pengajuan_jkk_plkk
  (
    p_kode_pengajuan    varchar2,
    p_user              varchar2,
    p_status_kelayakan  out varchar2, -- T = Tidak Layak Klaim, Y = Layak Klaim
    p_kode_pesan        out varchar2,
    p_pesan_tidak_layak out varchar2,
    p_sukses            out varchar2,
    p_mess              out varchar2
  ) as
  
    v_ada_pengajuan   number := 0;
    
    v_kode_pesan            varchar2(30):= '';
    v_pesan_tidak_layak     varchar2(4000):= '';
    v_status_kelayakan      varchar2(1):= 'T';
    
    v1 nsp.nsp_klaim@to_ec%rowtype; 
    v_kode_klaim            varchar2(30):= '';
  
    v_error     integer:=0;
    v_errmess   varchar2(4000):='';
    
    v_kode_jam_kecelakaan varchar2(30):= '';
  
  begin
    -- start cek data pengajuan status DRAFT
    select count(*) into v_ada_pengajuan
    from NSP.NSP_KLAIM@to_ec
    where nvl(STATUS_BATAL,'X') = 'T'
    and KODE_PENGAJUAN = p_kode_pengajuan 
    and KODE_STATUS_PENGAJUAN = 'ST01' -- ST01	DRAFT
    ;
    
    if v_ada_pengajuan > 0 then
      begin
        select a.* into v1
        from NSP.NSP_KLAIM@to_ec a
        where nvl(a.STATUS_BATAL,'X') = 'T'
        and a.KODE_PENGAJUAN = p_kode_pengajuan 
        and KODE_STATUS_PENGAJUAN = 'ST01'
        ; 
      exception when others then
        v_kode_pesan := 'JKKA999';
        -- LAYAK JIKA TIDAK TERDAPAT KESALAHAN PADA TRANSAKSI
        v_error := nvl(v_error,0)+1;
        v_errmess := 'Select data dari tabel NSP.NSP_KLAIM gagal. Error: '||sqlerrm;
      end;
    
      if(v1.KODE_JAM_KECELAKAAN is not null) then
      begin
        select kode_fix
        into
        v_kode_jam_kecelakaan
        from (
        select 
        case 
        when jam_cek >= jam_awal and jam_cek <= jam_akhir then kode
        else 'T' end as kode_fix
        from (
        (select 
        kode,
        to_number((regexp_replace(replace(v1.KODE_JAM_KECELAKAAN,'00:00','24:00'),' |:',''))) jam_cek,
        to_number(substr(regexp_replace(keterangan,' |:',''),1,4)) jam_awal,
        to_number(substr(regexp_replace(keterangan,' |:',''),-4)) jam_akhir
                                       from sijstk.ms_lookup
                                      where     tipe = 'KLMJAMKERJ'
                                      and aktif = 'Y')
        ) a
        )
        where kode_fix <> 'T'
        and rownum = 1; 
       exception when others then
        v_kode_pesan := 'JKKA999';
        v_error := nvl(v_error,0)+1;
        v_errmess := 'Select data dari tabel sijstk.ms_lookup gagal. Error: '||sqlerrm;
       end; 
       end if;
      
      -- start pengecekan apakah terdapat data mandatory yang kosong
      begin
        -- cek sebab klaim
        -- cek tanggal lapor
        -- cek tanggal kecelakaan
        -- cek jam kecelakaan
        -- cek jenis kasus
        -- cek lokasi kecelakaan
        -- cek tempat kecelakaan
        -- cek kronologi kecelakaan
        -- cek nama PLKK
        -- cek email TK
        -- cek no hp TK
        -- cek nama petugas perusahaan
        -- cek jabatan petugas perusahaan
        -- cek no hp petugas perusahaan
        -- cek email petugas perusahaan
        if v1.KODE_SEBAB_KLAIM = 'SKK01' then
          if v1.KODE_SEBAB_KLAIM is null 
            or v1.TGL_LAPOR is null
            or v1.TGL_KECELAKAAN is null
            or v1.KODE_JAM_KECELAKAAN is null
            or v1.KODE_JENIS_KASUS is null
            or v1.KODE_LOKASI_KECELAKAAN is null
            or v1.NAMA_TEMPAT_KECELAKAAN is null
            or v1.KRONOLOGI_KECELAKAAN is null
            or v1.KODE_TINDAKAN_BAHAYA is null
            or v1.KODE_KONDISI_BAHAYA is null
            or v1.KODE_CORAK is null
            or v1.KODE_SUMBER_CEDERA is null
            or v1.KODE_AKIBAT_DIDERITA is null --Problem ID 1638
            or v1.TRANSPORTASI is null --Problem ID 1638
            or v1.KODE_BAGIAN_SAKIT is null --Problem ID 1638
            or v1.KODE_PPK is null
            or v1.EMAIL is null 
            or v1.NOMOR_HP is null
            or v1.NAMA_PIC_PERUSAHAAN is null
            or v1.JABATAN_PIC_PERUSAHAAN is null
            or v1.NO_HP_PIC_PERUSAHAAN is null
            or v1.EMAIL_PIC_PERUSAHAAN is null then
              v_error := nvl(v_error,0)+1;
              v_errmess := 'Data informasi kecelakaan, tenaga kerja dan informasi pihak pelapor dari perusahaan harus diisi semua, silakan melengkapi data inputan terlebih dahulu dan klik tombol Simpan.';  
          end if;
        else
          if v1.KODE_SEBAB_KLAIM is null 
            or v1.TGL_LAPOR is null
            or v1.TGL_KECELAKAAN is null
            or v1.KODE_JENIS_KASUS is null
            or v1.KODE_PPK is null
            or v1.EMAIL is null 
            or v1.NOMOR_HP is null
            or v1.NAMA_PIC_PERUSAHAAN is null
            or v1.JABATAN_PIC_PERUSAHAAN is null
            or v1.NO_HP_PIC_PERUSAHAAN is null
            or v1.EMAIL_PIC_PERUSAHAAN is null then
              v_error := nvl(v_error,0)+1;
              v_errmess := 'Data informasi kecelakaan, tenaga kerja dan informasi pihak pelapor dari perusahaan harus diisi semua, silakan melengkapi data inputan terlebih dahulu dan klik tombol Simpan.';  
          end if;
        end if;
        
      exception when others then
        v_kode_pesan := 'JKKA999';
        -- LAYAK JIKA TIDAK TERDAPAT KESALAHAN PADA TRANSAKSI
        v_error := nvl(v_error,0)+1;
        v_errmess := 'Pengecekan data mandatory gagal. Error: '||sqlerrm;
      end;
      -- end pengecekan apakah terdapat data mandatory yang kosong
      
      if nvl(v_error,0)=0 then
       
        
        if nvl(v_error,0)=0 then
          begin
            -- update data pengajuan dari DRAFT menjadi DIAJUKAN
            update NSP.NSP_KLAIM@to_ec
            set KODE_STATUS_PENGAJUAN = 'ST03', -- ST03 LAPOR JKK TAHAP I
            STATUS_SUBMIT = 'Y',
            TGL_SUBMIT = sysdate,
            PETUGAS_SUBMIT = p_user,
            TGL_UBAH = sysdate,
            PETUGAS_UBAH = p_user
            where KODE_PENGAJUAN = p_kode_pengajuan
            and KODE_STATUS_PENGAJUAN = 'ST01' -- ST01	DRAFT
            ;
          exception when others then
            v_kode_pesan := 'JKKA999';
            v_error := nvl(v_error,0)+1;
            v_errmess := 'Submit data klaim gagal. Error: '||sqlerrm;
          end;
          
          -- start integrasi data ke PLKK
                   
          if nvl(v_error,0)=0 then
                      
            begin
             
               update tc.jm_tc_tahap1 set 
               kondisi_bahaya = (select keterangan from sijstk.ms_lookup where tipe='KLMKONDBHY' and nvl(aktif,'T')='Y' and kode=v1.kode_kondisi_bahaya),--P_KODE_KONDISI_BAHAYA,
               tindakan_bahaya = (select keterangan from sijstk.ms_lookup where tipe='KLMTINDBHY' and nvl(aktif,'T')='Y' and kode=v1.kode_tindakan_bahaya),--P_KODE_TINDAKAN_BAHAYA,
               corak = (select keterangan from sijstk.ms_lookup where tipe='KLMCORAK' and nvl(aktif,'T')='Y' and kode=v1.kode_corak),--P_KODE_CORAK,
               sumber_cidera =(select keterangan from sijstk.ms_lookup where tipe='KLMSMBRCDR' and nvl(aktif,'T')='Y' and kode=v1.kode_sumber_cedera),--P_KODE_SUMBER_CEDERA,
               bagian_yang_sakit =(select keterangan from sijstk.ms_lookup where tipe='KLMBGSAKIT' and nvl(aktif,'T')='Y' and kode=v1.kode_bagian_sakit),--P_KODE_BAGIAN_SAKIT,
               akibat =(select nama_akibat_diderita from pn.pn_kode_akibat_diderita where status_nonaktif='T' and kode_akibat_diderita=v1.kode_akibat_diderita),--P_KODE_AKIBAT_DIDERITA,
               status_submit = 'Y',
               st_app_pengajuan_kanal = 'Y',
               uraian_kronologi = v1.kronologi_kecelakaan,
               id_pelaporan_asal = p_kode_pengajuan,
               /*jenis_pekerjaan = v1.jenis_pekerjaan,*/ /*sudah diupdate oleh update detil klaim jkk ws sipp, update by davin hari jumat, 06-06-2025*/
               tgl_ubah = sysdate,
               petugas_ubah =  substr(p_user, 1, 30)
               where no_agenda = v1.kode_klaim;
               
               --28042025
                update tc.jm_tc_tahap1
                set tgl_3kkpak1 = sysdate,
                st_3kkpak1 = 'Y'
                where no_agenda = v1.kode_klaim
                and tgl_3kkpak1 is null;
                
                update tc.tc_klaim_dugaan_kk_pak
                set tgl_submit_3kk1 = sysdate
                where no_agenda = v1.kode_klaim
                and tgl_submit_3kk1 is null;
                         
                             
            exception when others then
              v_kode_pesan := 'JKKA999';
              v_error := nvl(v_error,0)+1;
              v_errmess := 'Integrasi data klaim JKK laporan tahap I ke PLKK untuk laporan tahap I gagal. Error: '||sqlerrm;
            end;
            
             if nvl(v_error,0)=0 then 
              begin
              
              INSERT INTO TC.TC_KLAIM_DOKUMEN (NO_AGENDA,
                                 TAHAP,
                                 KODE_DOKUMEN,
                                 MIME_TYPE,
                                 PATH_URL,
                                 FLAG_MANDATORY,
                                 FLAG_UPLOAD,
                                 TGL_UPLOAD,
                                 KETERANGAN,
                                 TGL_REKAM,
                                 PETUGAS_REKAM,
                                 TGL_UBAH,
                                 PETUGAS_UBAH)
            SELECT v1.KODE_KLAIM,
                   '1',
                   KODE_DOKUMEN,
                   MIME_TYPE,
                   PATH_URL,
                   FLAG_MANDATORY,
                   FLAG_UPLOAD,
                   TGL_UPLOAD,
                   KETERANGAN,
                   TGL_REKAM,
                   PETUGAS_REKAM,
                   TGL_UBAH,
                   PETUGAS_UBAH FROM NSP.NSP_KLAIM_DOKUMEN@TO_EC WHERE KODE_DOKUMEN IN 
                    ('D401','D111','D107') AND KODE_PENGAJUAN = p_kode_pengajuan;
              
               exception when others then
                v_kode_pesan := 'JKKA999';
                v_error := nvl(v_error,0)+1;
                v_errmess := 'Integrasi data klaim JKK laporan tahap I ke PLKK untuk dokumen gagal. Error: '||sqlerrm;
              end;
            end if;  
            
            if nvl(v_error,0)=0 then
            
            begin 
            
            MERGE INTO PN.PN_KLAIM_DOKUMEN t1
                USING
                (
                SELECT 
                v1.KODE_KLAIM KODE_KLAIM,
                CASE WHEN KODE_DOKUMEN = 'D400' THEN 'D130'
                WHEN KODE_DOKUMEN = 'D401' THEN 'D083'
                ELSE KODE_DOKUMEN END KODE_DOKUMEN,
                MIME_TYPE,
                PATH_URL,
                (SELECT NAMA_DOKUMEN FROM NSP.NSP_KLAIM_KODE_DOKUMEN@TO_EC WHERE KODE_DOKUMEN = A.KODE_DOKUMEN) NAMA_DOKUMEN,
                FLAG_UPLOAD,
                TGL_UPLOAD,
                TGL_UBAH,
                'SISTEM' PETUGAS_UBAH       
                FROM NSP.NSP_KLAIM_DOKUMEN@TO_EC A WHERE KODE_PENGAJUAN = p_kode_pengajuan AND PATH_URL IS NOT NULL
                )t2
                ON(t1.KODE_KLAIM = t2.KODE_KLAIM AND t1.KODE_DOKUMEN = t2.KODE_DOKUMEN)
                WHEN MATCHED THEN UPDATE SET
                t1.MIME_TYPE = t2.MIME_TYPE,
                t1.URL = t2.PATH_URL,
                t1.NAMA_FILE = T2.NAMA_DOKUMEN,
                t1.STATUS_DISERAHKAN = t2.FLAG_UPLOAD,
                t1.TGL_DISERAHKAN = SYSDATE,
                t1.TGL_UBAH = t2.TGL_UBAH,
                t1.PETUGAS_UBAH = t2.PETUGAS_UBAH;

              exception when others then
                v_kode_pesan := 'JKKA999';
                v_error := nvl(v_error,0)+1;
                v_errmess := 'Integrasi data klaim JKK laporan tahap I ke SMILE gagal. Error: '||sqlerrm;
              end;
            end if;
              
            
            if nvl(v_error,0)=0 then
            
              begin
                  update pn.pn_klaim set 
                  kode_kondisi_bahaya = v1.kode_kondisi_bahaya,
                  kode_tindakan_bahaya = v1.kode_tindakan_bahaya,
                  kode_corak = v1.kode_corak,
                  kode_sumber_cedera = v1.kode_sumber_cedera,
                  kode_bagian_sakit = v1.kode_bagian_sakit,
                  kode_akibat_diderita = v1.kode_akibat_diderita
--                  kode_jam_kecelakaan = v_kode_jam_kecelakaan
                  where kode_klaim = v1.kode_klaim;
                        
              exception when others then
                v_kode_pesan := 'JKKA999';
                v_error := nvl(v_error,0)+1;
                v_errmess := 'Integrasi data klaim JKK laporan tahap I ke SMILE gagal. Error: '||sqlerrm;
              end;
            end if;
            
             if nvl(v_error,0)=0 then
            
              begin
                  update nsp.nsp_klaim@to_ec set 
                  kode_status_pengajuan = 'ST03',
                  status_submit = 'Y'
                  where kode_klaim = v1.kode_klaim;
                        
              exception when others then
                v_kode_pesan := 'JKKA999';
                v_error := nvl(v_error,0)+1;
                v_errmess := 'Integrasi data klaim JKK laporan tahap I ke SIPP gagal. Error: '||sqlerrm;
              end;
            end if;

          end if;
          -- end integrasi data ke PLKK
        end if;
      end if;
    else
      v_kode_pesan := 'JKKA001';
      -- LAYAK JIKA DATA PENGAJUAN VALID
      v_error := nvl(v_error,0)+1;
      v_errmess := 'Layak jika data pengajuan valid.';
    end if;
    
    p_status_kelayakan := v_status_kelayakan;
    
    if nvl(v_error,0)=0 then
      p_sukses := '1'; 
      p_mess := 'Ok'; 
      commit;
      
    else
      p_sukses := '-1';
      p_mess := v_errmess;
      rollback;
    end if;
  end x_submit_pengajuan_jkk_plkk;
  