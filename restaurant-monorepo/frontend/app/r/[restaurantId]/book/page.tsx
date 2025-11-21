"use client";

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { useParams } from 'next/navigation';

export default function PublicBookingPage() {
  const params = useParams();
  const restaurantId = params.restaurantId as string;
  
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [guests, setGuests] = useState('2');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Combine date and time for StartTime
      const startTime = new Date(`${date}T${time}:00`).toISOString();
      
      await api.post(`/public/restaurants/${restaurantId}/bookings`, {
        CustomerName: name,
        CustomerEmail: email,
        StartTime: startTime,
        NumberOfGuests: parseInt(guests)
      });
      setSuccess(true);
    } catch (e) {
      console.error(e);
      alert('Booking failed');
    }
  };

  if (success) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-green-50">
        <Card className="w-[400px] text-center">
          <CardHeader>
            <CardTitle className="text-green-600">Booking Confirmed!</CardTitle>
          </CardHeader>
          <CardContent>
            <p>Thank you, {name}. Your table for {guests} is reserved.</p>
            <Button className="mt-4" onClick={() => setSuccess(false)}>Make Another Booking</Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle>Book a Table</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label>Name</Label>
              <Input value={name} onChange={e => setName(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label>Email</Label>
              <Input type="email" value={email} onChange={e => setEmail(e.target.value)} required />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label>Date</Label>
                <Input type="date" value={date} onChange={e => setDate(e.target.value)} required />
              </div>
              <div className="space-y-2">
                <Label>Time</Label>
                <Input type="time" value={time} onChange={e => setTime(e.target.value)} required />
              </div>
            </div>
            <div className="space-y-2">
              <Label>Number of Guests</Label>
              <Input type="number" min="1" value={guests} onChange={e => setGuests(e.target.value)} required />
            </div>
            <Button type="submit" className="w-full">Confirm Booking</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
