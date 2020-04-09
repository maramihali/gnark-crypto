// Code generated by internal/pairing DO NOT EDIT
package bls377

import (
	"testing"

	"github.com/consensys/gurvy/bls377/fp"
	"github.com/consensys/gurvy/bls377/fr"
)

func TestPairingLineEval(t *testing.T) {

	G := G2Jac{}
	G.X.SetString("11467063222684898633036104763692544506257812867640109164430855414494851760297509943081481005947955008078272733624",
		"153924906120314059329163510034379429156688480181182668999642334674073859906019623717844462092443710331558842221198")
	G.Y.SetString("217426664443013466493849511677243421913435679616098405782168799962712362374085608530270502677771125796970144049342",
		"220113305559851867470055261956775835250492241909876276448085325823827669499391027597256026508256704101389743638320")
	G.Z.SetString("1",
		"0")

	H := G2Jac{}
	H.X.SetString("38348804106969641131654336618231918247608720362924380120333996440589719997236048709530218561145001033408367199467",
		"208837221672103828632878568310047865523715993428626260492233587961023171407529159232705047544612759994485307437530")
	H.Y.SetString("219129261975485221488302932474367447253380009436652290437731529751224807932621384667224625634955419310221362804739",
		"62857965187173987050461294586432573826521562230975685098398439555961148392353952895313161290735015726193379258321")
	H.Z.SetString("1",
		"0")

	var a, b, c fp.Element
	a.SetString("219129261975485221488302932474367447253380009436652290437731529751224807932621384667224625634955419310221362804739")
	b.SetString("62857965187173987050461294586432573826521562230975685098398439555961148392353952895313161290735015726193379258321")
	c.SetString("1")
	P := G1Jac{}
	P.X = a
	P.Y = b
	P.Z = c

	var Paff G1Affine
	P.ToAffineFromJac(&Paff)

	el := &e12{}
	lRes := &lineEvalRes{}
	lineEvalJac(G, H, &Paff, lRes)
	el.C0.B1 = lRes.r0
	el.C1.B1 = lRes.r1
	el.C1.B2 = lRes.r2

	// el.FromMont()

	expected := "0+0*u+(220291599185938038585565774521033812062947190299680306664648725201730830885666933651848261361463591330567860207241+232134458700276476669584229661634543747068594368664068937164975724095736595288995356706959089579876199020312643174*u)*v+(0+0*u)*v**2+(0+0*u+(74241662856820718491669277383162555524896537826488558937227282983357670568906847284642533051528779250776935382660+9787836945036920457066634104342154603142239983688979247440278426242314457905122599227144555989168817796094251258*u)*v+(85129589817387660717039592198118788807152207633847410148299763250229022303850156734979397272700502238285752744807+245761211327131018855579902758747359135620549826797077633679496719449586668701082009536667506317412690997533857875*u)*v**2)*w"

	if expected != el.String() {
		t.Fatal("expected", expected, "got", el.String())
	}
}

func TestMagicPairing(t *testing.T) {
	var r1, r2 e12

	r1.SetRandom()
	r2.SetRandom()

	t.Log(r1)
	t.Log(r2)

	curve := BLS377()

	res1 := curve.FinalExponentiation(&r1)
	res2 := curve.FinalExponentiation(&r2)

	if res1.Equal(&res2) {
		t.Fatal("TestMagicPairing failed")
	}
}

func TestComputePairing(t *testing.T) {

	curve := BLS377()

	G := curve.g2Gen.Clone()
	P := curve.g1Gen
	sG := G.Clone()
	sP := P.Clone()

	var Gaff, sGaff G2Affine
	var Paff, sPaff G1Affine

	// checking bilinearity

	// check 1
	scalar := fr.Element{123}
	sG.ScalarMul(curve, G, scalar)
	sP.ScalarMul(curve, sP, scalar)

	var mRes1, mRes2, mRes3 e12

	P.ToAffineFromJac(&Paff)
	sP.ToAffineFromJac(&sPaff)
	G.ToAffineFromJac(&Gaff)
	sG.ToAffineFromJac(&sGaff)

	res1 := curve.FinalExponentiation(curve.MillerLoop(Paff, sGaff, &mRes1))
	res2 := curve.FinalExponentiation(curve.MillerLoop(sPaff, Gaff, &mRes2))

	if !res1.Equal(&res2) {
		t.Fatal("pairing failed")
	}

	// check 2
	s1G := G.Clone()
	s2G := G.Clone()
	s3G := G.Clone()
	s1 := fr.Element{29372983}
	s2 := fr.Element{209302420904}
	var s3 fr.Element
	s3.Add(&s1, &s2)
	s1G.ScalarMul(curve, s1G, s1)
	s2G.ScalarMul(curve, s2G, s2)
	s3G.ScalarMul(curve, s3G, s3)

	var s1Gaff, s2Gaff, s3Gaff G2Affine
	s1G.ToAffineFromJac(&s1Gaff)
	s2G.ToAffineFromJac(&s2Gaff)
	s3G.ToAffineFromJac(&s3Gaff)

	rs1 := curve.FinalExponentiation(curve.MillerLoop(Paff, s1Gaff, &mRes1))
	rs2 := curve.FinalExponentiation(curve.MillerLoop(Paff, s2Gaff, &mRes2))
	rs3 := curve.FinalExponentiation(curve.MillerLoop(Paff, s3Gaff, &mRes3))
	rs1.Mul(&rs2, &rs1)
	if !rs3.Equal(&rs1) {
		t.Fatal("pairing failed2")
	}

}

func TestFinalExponentiation(t *testing.T) {

	var a e12

	// curve initialization
	curve := BLS377()

	a.SetString(
		"1382424129690940106527336948935335363935127549146605398842626667204683483408227749",
		"0121296909401065273369489353353639351275491466053988426266672046834834082277499690",
		"7336948129690940106527336948935335363935127549146605398842626667204683483408227749",
		"6393512129690940106527336948935335363935127549146605398842626667204683483408227749",
		"2581296909401065273369489353353639351275491466053988426266672046834834082277496644",
		"5331296909401065273369489353353639351275491466053988426266672046834834082277495363",
		"1296909401065273369489353353639351275491466053988426266672046834834082277491382424",
		"0129612969094010652733694893533536393512754914660539884262666720468348340822774990",
		"7336948129690940106527336948935335363935127549146605398842626667204683483408227749",
		"6393129690940106527336948935335363935127549146605398842626667204683483408227749512",
		"2586641296909401065273369489353353639351275491466053988426266672046834834082277494",
		"5312969094010652733694893533536393512754914660539884262666720468348340822774935363")

	got := curve.FinalExponentiation(&a)
	// got.FromMont()
	expected := "144194728711271115200960338354399242675728583773229950397819045476446438084204609471396642031402812685868991569321+211655715903640692993096217798002440469463348812262035874643014089609279211697384357585411910878072499680258137966*u+(142847066272809567418100082074671069509232241577000969140908287562600131374001244027677233030439295603512969014098+195798360278497370480514919830707791165439804107831313818538392791182098078978554624847462732535644384771463893749*u)*v+(86535834900386741976161849459901164117916658028045848198896601307003900317367651452650171453989007615756269982055+224852590401918627972000608877264354401740654844913628932015173271867650456466791238805911444161155014508555986302*u)*v**2+(84753048733212221347206140939947896362463706184166407933992007029321604869923850544953196752734135050061759565628+65455822154707898400320891026725436969464960341136383846714152247599725536100703394287883773613448858651461616688*u+(166022890778297035042874386744221331185425031233249433341140091448935862547653045382140453272792874748437049148206+112097397635628424819887844002083726992692863111606459600443975137441512833688639214196992677192944546581733745228*u)*v+(154270191003785238867416493510207961994765591028485520151537180047545003150612330145604984518433309900477076500369+71053829709457555512251373435255765851107607797752394326222141063282430248152312493662780307517882687885460943678*u)*v**2)*w"

	if expected != got.String() {
		t.Fatal("expected", expected, "got", got.String())
	}

}

//--------------------//
//     benches		  //
//--------------------//

func BenchmarkLineEval(b *testing.B) {

	curve := BLS377()

	H := G2Jac{}
	H.ScalarMul(curve, &curve.g2Gen, fr.Element{1213})

	lRes := &lineEvalRes{}
	var g1GenAff G1Affine
	curve.g1Gen.ToAffineFromJac(&g1GenAff)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lineEvalJac(curve.g2Gen, H, &g1GenAff, lRes)
	}

}

func BenchmarkPairing(b *testing.B) {

	curve := BLS377()

	var mRes e12

	var g1GenAff G1Affine
	var g2GenAff G2Affine

	curve.g1Gen.ToAffineFromJac(&g1GenAff)
	curve.g2Gen.ToAffineFromJac(&g2GenAff)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.FinalExponentiation(curve.MillerLoop(g1GenAff, g2GenAff, &mRes))
	}
}

func BenchmarkFinalExponentiation(b *testing.B) {

	var a e12

	curve := BLS377()

	a.SetString(
		"1382424129690940106527336948935335363935127549146605398842626667204683483408227749",
		"0121296909401065273369489353353639351275491466053988426266672046834834082277499690",
		"7336948129690940106527336948935335363935127549146605398842626667204683483408227749",
		"6393512129690940106527336948935335363935127549146605398842626667204683483408227749",
		"2581296909401065273369489353353639351275491466053988426266672046834834082277496644",
		"5331296909401065273369489353353639351275491466053988426266672046834834082277495363",
		"1296909401065273369489353353639351275491466053988426266672046834834082277491382424",
		"0129612969094010652733694893533536393512754914660539884262666720468348340822774990",
		"7336948129690940106527336948935335363935127549146605398842626667204683483408227749",
		"6393129690940106527336948935335363935127549146605398842626667204683483408227749512",
		"2586641296909401065273369489353353639351275491466053988426266672046834834082277494",
		"5312969094010652733694893533536393512754914660539884262666720468348340822774935363")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		curve.FinalExponentiation(&a)
	}

}