package dusk

import (
	"time"
)

type EquatorialCoordinate struct {
	ra  float64
	dec float64
}

type EclipticCoordinate struct {
	λ float64
	β float64
	Δ float64
}

type HorizontalCoordinate struct {
	/*
		altitude (a) or elevation
	*/
	a float64
	/*
		azimuth (A) or elevation
	*/
	A float64
}

/*
	ConvertEclipticCoordinateToEquatorial()

	@param datetime - the datetime of the observer (in UTC)
	@param geocentric ecliptic coordinate of type EclipticCoordinate { λ, β, Λ }
	@returns the converted equatorial coordinate { ra, dec }
	@see eq13.3 & eq13.4 p.93 of Meeus, Jean. 1991. Astronomical algorithms. Richmond, Va: Willmann-Bell.
*/
func ConvertEclipticCoordinateToEquatorial(datetime time.Time, ec EclipticCoordinate) EquatorialCoordinate {
	var J = GetCurrentJulianCenturyRelativeToJ2000(datetime)

	var L float64 = GetSolarMeanLongitude(J)

	var l float64 = GetLunarMeanLongitude(J)

	var Ω float64 = GetLunarLongitudeOfTheAscendingNode(J)

	var ε float64 = GetMeanObliquityOfTheEcliptic(J) + GetNutationInObliquityOfTheEcliptic(L, l, Ω)

	var λ = ec.λ

	var β = ec.β

	var α = atanx((sinx(λ)*cosx(ε) - tanx(β)*sinx(ε)) / cosx(λ))

	var δ = asinx(sinx(β)*cosx(ε) + cosx(β)*sinx(ε)*sinx(λ))

	if α < 0 {
		α += 180
	}

	return EquatorialCoordinate{
		ra:  α,
		dec: δ,
	}
}

/*
	ConvertEquatorialCoordinateToHorizontal()

	@param datetime - the datetime of the observer (in UTC)
	@param longitude - is the longitude (west is negative, east is positive) in degrees of some observer on Earth
	@param latitude - is the latitude (south is negative, north is positive) in degrees of some observer on Earth
	@param equatorial coordinate of type EquatorialCoordiate { ra, dec }
	@returns the equivalent horizontal coordinate for the given observers position
	@see eq13.5 and eq.6 p.93 of Meeus, Jean. 1991. Astronomical algorithms. Richmond, Va: Willmann-Bell.
*/
func ConvertEquatorialCoordinateToHorizontal(datetime time.Time, longitude float64, latitude float64, eq EquatorialCoordinate) HorizontalCoordinate {
	var LST float64 = GetLocalSiderealTime(datetime, longitude)

	var ra float64 = GetHourAngle(eq.ra, LST)

	var dec float64 = eq.dec

	var alt = asinx(sinx(dec)*sinx(latitude) + cosx(dec)*cosx(latitude)*cosx(ra))

	var az = acosx((sinx(dec) - sinx(alt)*sinx(latitude)) / (cosx(alt) * cosx(latitude)))

	return HorizontalCoordinate{
		a: alt,
		A: az,
	}
}
