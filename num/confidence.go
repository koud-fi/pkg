package num

import "math"

func ConfidenceVote(upVotes, downVotes int) float64 {
	if upVotes < 0 {
		upVotes = 0
	}
	if downVotes < 0 {
		downVotes = 0
	}
	n := upVotes + downVotes
	if n == 0 {
		return 0.0
	}
	p := float64(upVotes) / float64(n)
	return Confidence(p, n)
}

func ConfidenceScore(score, min, max float64, voteCount int) float64 {
	if min == max {
		return 0.0
	} else if min > max {
		min, max = max, min
	}
	nf := float64(voteCount)
	if score > nf*max {
		score = nf * max
	}
	if score < nf*min {
		score = nf * min
	}
	p := (score - min*nf) / ((max - min) * nf)
	return Confidence(p, voteCount)
}

func Confidence(p float64, n int) float64 {
	if n <= 0 {
		return 0.0
	}
	if p < 0.0 {
		p = 0.0
	} else if p > 1.0 {
		p = 1.0
	}
	var (
		z  = 1.281551565545
		nf = float64(n)

		left  = p + 1/(2*nf)*z*z
		right = z * math.Sqrt(p*(1-p)/nf+z*z/(4*nf*nf))
		under = 1 + 1/nf*z*z
	)
	return (left - right) / under
}
